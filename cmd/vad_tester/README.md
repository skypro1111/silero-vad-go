# Silero VAD Tester

A command-line tool for testing and experimenting with the Silero Voice Activity Detection (VAD) model.

## Building

```bash
go build -o vad_tester
```

## Usage

```bash
./vad_tester [options]
```

### Command Line Options

- `-model` - Path to the Silero VAD ONNX model file (default: `testfiles/silero_vad.onnx`)
- `-audio` - Path to the PCM audio file (default: `testfiles/samples.pcm`)
- `-sr` - Sample rate, either 8000 or 16000 Hz (default: 16000)
- `-threshold` - Speech detection probability threshold (default: 0.5)
- `-neg-threshold` - Silence detection probability threshold (default: 0.0 = auto)
- `-min-silence` - Minimum silence duration in milliseconds (default: 500)
- `-min-speech` - Minimum speech duration in milliseconds (default: 250)
- `-speech-pad` - Speech segments padding in milliseconds (default: 30)
- `-verbose` - Enable verbose output (default: false)

### Parameter Tuning

For comparing different parameter combinations, you can use the provided scripts:

```bash
# Test with first audio file
./compare_params.sh

# Test with second audio file
./compare_params2.sh
```

## Understanding Parameters

### Speech Detection Threshold

The probability threshold above which audio is classified as speech. Range: (0, 1).
- Higher values (e.g., 0.7) = more conservative, less false positives
- Lower values (e.g., 0.3) = more sensitive, might catch quieter speech
- Default (0.5) is a good balance for most cases

### Silence Detection Threshold

The probability threshold below which audio is classified as silence. Range: [0, 1).
- 0.0 = automatically set to (threshold - 0.15)
- Higher values = faster switching to silence
- Lower values = more persistent speech detection
- Must be less than speech threshold

### Minimum Silence Duration

The minimum duration of silence required to split speech segments.
- Higher values (e.g., 1000ms) = fewer splits, longer segments
- Lower values (e.g., 200ms) = more splits, shorter segments
- Default (500ms) works well for normal speech

### Minimum Speech Duration

The minimum duration required for a segment to be considered valid speech.
- Higher values = filter out short noises
- Lower values = catch brief utterances
- Default (250ms) catches most meaningful speech while filtering noise

### Speech Padding

Padding added to the beginning and end of speech segments.
- Higher values = less aggressive cutting, might include more silence
- Lower values = tighter segments, might cut speech edges
- Default (30ms) provides good balance

## Examples

1. Basic usage with default parameters:
```bash
./vad_tester -model path/to/model.onnx -audio path/to/audio.pcm
```

2. More sensitive speech detection:
```bash
./vad_tester -threshold 0.3 -neg-threshold 0.1 -min-speech 150
```

3. Conservative detection with longer segments:
```bash
./vad_tester -threshold 0.7 -min-silence 1000 -min-speech 500
```

4. Fine-tuned for quick responses:
```bash
./vad_tester -threshold 0.5 -neg-threshold 0.4 -min-silence 200 -min-speech 200 -speech-pad 10
```

## Output Format

The tool provides:
1. Information about the audio file
2. Detection parameters used
3. List of detected speech segments with timestamps
4. Total duration statistics
5. Real-time factor (RTF) performance metric

Example output:
```
Detected speech segments:
------------------------
1. 1.06 - 1.63 (0.57 sec)
2. 2.88 - 3.23 (0.35 sec)
3. 4.45 - 4.89 (0.44 sec)

Total audio duration: 5.00 sec
Total speech duration: 1.36 sec (27.2%)
``` 