import asyncio
import websockets
import json
import base64
import logging
from typing import Dict
import uuid
import os
import aiohttp
from datetime import datetime

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.StreamHandler(),
        logging.FileHandler('interview_service.log')
    ]
)
logger = logging.getLogger(__name__)

class InterviewManager:
    def __init__(self):
        self.sessions: Dict[str, dict] = {}
        self.deepgram_api_key = os.getenv('DEEPGRAM_API_KEY')
        self.ai_service_url = os.getenv('AI_SERVICE_URL', 'http://ai-service:8082')

        if not self.deepgram_api_key:
            logger.error("DEEPGRAM_API_KEY not set. Speech services will not work properly.")

    def create_session(self, candidate_id: str, vacancy: dict, resume_text: str) -> str:
        session_id = str(uuid.uuid4())
        self.sessions[session_id] = {
            'candidate_id': candidate_id,
            'vacancy': vacancy,
            'resume_text': resume_text,
            'current_question': 0,
            'questions': [],
            'answers': [],
            'score': 0,
            'start_time': datetime.now().isoformat(),
            'status': 'waiting'
        }
        logger.info(f"Создана сессия интервью: {session_id}")
        return session_id

    async def generate_questions(self, session_id: str):
        session = self.sessions.get(session_id)
        if not session:
            return None

        try:
            async with aiohttp.ClientSession() as http_session:
                data = {
                    "vacancy": session['vacancy'],
                    "resume_text": session['resume_text']
                }

                async with http_session.post(
                    f"{self.ai_service_url}/generate-questions",
                    json=data,
                    timeout=aiohttp.ClientTimeout(total=30)
                ) as response:
                    if response.status == 200:
                        result = await response.json()
                        
                        questions = result.get('questions', [])
                        session['questions'] = questions
                        logger.info(f"Сгенерировано вопросов: {len(session['questions'])}")
                        return session['questions']
                    else:
                        logger.error(f"Ошибка генерации вопросов: {response.status}")
                        return None
        except Exception as e:
            logger.error(f"Исключение при генерации вопросов: {e}")
            return None

    async def analyze_answer(self, session_id: str, question: str, answer: str) -> float:
        session = self.sessions.get(session_id)
        if not session:
            return 0.0

        try:
            async with aiohttp.ClientSession() as http_session:
                data = {
                    "question": question,
                    "answer": answer,
                    "vacancy": session['vacancy']
                }

                async with http_session.post(
                    f"{self.ai_service_url}/analyze-answer",
                    json=data,
                    timeout=aiohttp.ClientTimeout(total=30)
                ) as response:
                    if response.status == 200:
                        result = await response.json()
                        analysis_str = result.get('result', '{}')  
                        try:
                            analysis = json.loads(analysis_str)  
                        except json.JSONDecodeError as e:
                            logger.error(f"Ошибка парсинга анализа JSON: {e}")
                            analysis = {}
                        score = analysis.get('score', 0.5)
                        logger.info(f"Анализ ответа: score={score}")
                        return score
                    else:
                        logger.error(f"Ошибка анализа ответа: {response.status}")
                        return 0.5
        except Exception as e:
            logger.error(f"Исключение при анализе ответа: {e}")
            return 0.5

    def get_next_question(self, session_id: str) -> str:
        session = self.sessions.get(session_id)
        if not session:
            return None

        if session['current_question'] < len(session['questions']):
            question = session['questions'][session['current_question']]
            logger.info(f"Вопрос {session['current_question'] + 1}: {question}")
            return question
        return None

    async def save_answer(self, session_id: str, answer: str, confidence: float):
        session = self.sessions.get(session_id)
        if not session:
            return False

        question = session['questions'][session['current_question']]

        ai_score = await self.analyze_answer(session_id, question, answer)

        final_score = confidence * ai_score

        session['answers'].append({
            'question': question,
            'answer': answer,
            'confidence': confidence,
            'ai_score': ai_score,
            'final_score': final_score,
            'timestamp': datetime.now().isoformat()
        })

        session['score'] += final_score
        session['current_question'] += 1

        logger.info(f"Сохранен ответ: confidence={confidence:.2f}, ai_score={ai_score:.2f}, final={final_score:.2f}")
        return True

    def get_results(self, session_id: str) -> dict:
        session = self.sessions.get(session_id)
        if not session:
            return None

        session['end_time'] = datetime.now().isoformat()
        session['duration'] = (datetime.fromisoformat(session['end_time']) -
                              datetime.fromisoformat(session['start_time'])).total_seconds()

        logger.info(f"Интервью завершено для сессии {session_id}")
        return session

    async def deepgram_text_to_speech(self, text: str) -> bytes:
        
        if not self.deepgram_api_key:
            return None

        url = 'https://api.deepgram.com/v1/speak?model=aura-2-thalia-en'
        headers = {
            'Authorization': f'Token {self.deepgram_api_key}',
            'Content-Type': 'text/plain'
        }

        try:
            async with aiohttp.ClientSession() as session:
                async with session.post(url, headers=headers, data=text, timeout=aiohttp.ClientTimeout(total=30)) as response:
                    if response.status == 200:
                        return await response.read()
                    else:
                        error_text = await response.text()
                        logger.error(f"Ошибка Deepgram TTS: {response.status}, {error_text}")
                        return None
        except Exception as e:
            logger.error(f"Исключение в Deepgram TTS: {e}")
            return None

    async def deepgram_speech_to_text(self, audio_data: bytes) -> str:
        
        if not self.deepgram_api_key:
            return "Ошибка: не настроен API ключ", 0.0

        params = {
            'detect_language': 'true',
            'model': 'nova-3',
            'smart_format': 'true',
            'punctuate': 'true'
        }

        url = 'https://api.deepgram.com/v1/listen'
        headers = {
            'Authorization': f'Token {self.deepgram_api_key}',
            'Content-Type': 'audio/webm'
        }

        try:
            async with aiohttp.ClientSession() as session:
                full_url = f"{url}?{'&'.join([f'{k}={v}' for k, v in params.items()])}"

                async with session.post(full_url, headers=headers, data=audio_data, timeout=aiohttp.ClientTimeout(total=30)) as response:
                    if response.status == 200:
                        result = await response.json()
                        transcript = result.get('results', {}).get('channels', [{}])[0].get('alternatives', [{}])[0]
                        transcribed_text = transcript.get('transcript', '')
                        confidence = transcript.get('confidence', 0.5)

                        language = result.get('results', {}).get('channels', [{}])[0].get('detected_language', 'unknown')
                        logger.info(f"Распознано: '{transcribed_text}', уверенность: {confidence}, язык: {language}")

                        return transcribed_text, confidence
                    else:
                        error_text = await response.text()
                        logger.error(f"Ошибка Deepgram STT: {response.status}, {error_text}")
                        return "Ошибка распознавания речи", 0.0
        except Exception as e:
            logger.error(f"Исключение в Deepgram STT: {e}")
            return "Ошибка соединения с сервисом распознавания", 0.0

interview_manager = InterviewManager()

async def handle_interview(websocket, path):
    client_ip = websocket.remote_address[0]
    logger.info(f"Новое подключение от {client_ip}")

    try:
        async for message in websocket:
            data = json.loads(message)
            logger.info(f"Получено сообщение типа: {data['type']}")

            if data['type'] == 'start_interview':
                
                session_id = interview_manager.create_session(
                    data['candidate_id'],
                    data['vacancy'],
                    data['resume_text']
                )

                questions = await interview_manager.generate_questions(session_id)
                if not questions:
                    await websocket.send(json.dumps({
                        'type': 'error',
                        'message': 'Не удалось сгенерировать вопросы для интервью'
                    }))
                    continue

                interview_manager.sessions[session_id]['status'] = 'active'

                first_question = interview_manager.get_next_question(session_id)

                await websocket.send(json.dumps({
                    'type': 'question_text',
                    'session_id': session_id,
                    'question': first_question,
                    'question_number': 1,
                    'total_questions': len(questions)
                }))

                question_audio = await interview_manager.deepgram_text_to_speech(first_question)
                if question_audio:
                    audio_base64 = base64.b64encode(question_audio).decode('utf-8')
                    await websocket.send(json.dumps({
                        'type': 'question_audio',
                        'session_id': session_id,
                        'question_audio': audio_base64
                    }))

                await websocket.send(json.dumps({
                    'type': 'ready_to_record',
                    'session_id': session_id
                }))

                question_audio = await interview_manager.deepgram_text_to_speech(first_question)
                if question_audio:
                    audio_base64 = base64.b64encode(question_audio).decode('utf-8')
                    await websocket.send(json.dumps({
                        'type': 'question_audio',
                        'session_id': session_id,
                        'question_audio': audio_base64
                    }))

            elif data['type'] == 'audio_response':
                
                session_id = data['session_id']
                logger.info(f"Обработка аудио ответа для сессии {session_id}")

                try:
                    await websocket.send(json.dumps({
                        'type': 'processing_started',
                        'session_id': session_id
                    }))

                    audio_data = base64.b64decode(data['audio'])
                    transcribed_text, confidence = await interview_manager.deepgram_speech_to_text(audio_data)

                    await interview_manager.save_answer(session_id, transcribed_text, confidence)

                    next_question = interview_manager.get_next_question(session_id)

                    if next_question:
                        
                        await websocket.send(json.dumps({
                            'type': 'question_text',
                            'session_id': session_id,
                            'question': next_question,
                            'question_number': interview_manager.sessions[session_id]['current_question'] + 1,
                            'total_questions': len(interview_manager.sessions[session_id]['questions'])
                        }))

                        question_audio = await interview_manager.deepgram_text_to_speech(next_question)
                        if question_audio:
                            audio_base64 = base64.b64encode(question_audio).decode('utf-8')
                            await websocket.send(json.dumps({
                                'type': 'question_audio',
                                'session_id': session_id,
                                'question_audio': audio_base64
                            }))
                    else:
                        
                        results = interview_manager.get_results(session_id)
                        await websocket.send(json.dumps({
                            'type': 'interview_completed',
                            'session_id': session_id,
                            'score': results['score'],
                            'answers': results['answers'],
                            'total_score': round(results['score'] * 25, 1),
                            'duration': round(results['duration'], 1)
                        }))

                except Exception as e:
                    logger.error(f"Ошибка обработки аудио: {e}")
                    await websocket.send(json.dumps({
                        'type': 'error',
                        'message': 'Ошибка обработки аудио ответа'
                    }))

    except websockets.exceptions.ConnectionClosed as e:
        logger.info(f"Соединение закрыто: {e}")
    except Exception as e:
        logger.error(f"Ошибка в обработчике интервью: {e}")

async def main():
    server = await websockets.serve(handle_interview, "0.0.0.0", 8765)
    logger.info("WebSocket сервер запущен на порту 8765")

    await server.wait_closed()

if __name__ == "__main__":
    asyncio.run(main())