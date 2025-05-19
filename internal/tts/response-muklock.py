from gtts import gTTS
import sys

text = sys.argv[1]
tts = gTTS(text=text, lang='th')
tts.save("response-muklock.mp3")