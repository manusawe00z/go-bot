from gtts import gTTS
import sys
import os
import tempfile
import logging
from pathlib import Path

# Set up logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger('tts')

def text_to_speech(text, lang='en', output_file='tts.mp3'):
    """
    Convert text to speech using Google Text-to-Speech.
    
    Args:
        text (str): The text to convert to speech
        lang (str): The language code (default: 'en')
        output_file (str): Path to save the audio file
    
    Returns:
        bool: True if successful, False otherwise
    """
    try:
        logger.info(f"Converting text to speech: '{text}' (language: {lang})")
        tts = gTTS(text=text, lang=lang, slow=False)
        
        # Create a temporary file first to avoid issues with concurrent access
        with tempfile.NamedTemporaryFile(suffix='.mp3', delete=False) as temp_file:
            temp_path = temp_file.name
        
        # Save to temporary file
        tts.save(temp_path)
        
        # Move to final destination
        os.replace(temp_path, output_file)
        
        logger.info(f"Audio saved to {output_file}")
        return True
    except Exception as e:
        logger.error(f"Error in text_to_speech: {str(e)}")
        return False

if __name__ == "__main__":
    # Get text from command line arguments
    if len(sys.argv) < 2:
        logger.error("Usage: python tts.py <text> [language] [output_file]")
        sys.exit(1)
    
    text = sys.argv[1]
    
    # Get language from arguments or use default
    lang = 'th'
    if len(sys.argv) > 2:
        lang = sys.argv[2]
    
    # Get output file from arguments or use default
    output_file = 'tts.mp3'
    if len(sys.argv) > 3:
        output_file = sys.argv[3]
    
    # Ensure the directory exists
    output_dir = os.path.dirname(output_file)
    if output_dir and not os.path.exists(output_dir):
        os.makedirs(output_dir)
    
    # Convert text to speech
    success = text_to_speech(text, lang, output_file)
    sys.exit(0 if success else 1)