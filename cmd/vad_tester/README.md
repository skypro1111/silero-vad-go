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

### Basic Usage

1. Process an audio file with default parameters:
```bash
./vad_tester -model path/to/model.onnx -audio path/to/audio.pcm
```

2. Use with 8kHz audio:
```bash
./vad_tester -model path/to/model.onnx -audio path/to/8khz_audio.pcm -sr 8000
```

3. Enable verbose output for debugging:
```bash
./vad_tester -model path/to/model.onnx -audio path/to/audio.pcm -verbose
```

### Parameter Tuning Examples

1. More sensitive speech detection (good for quiet speech):
```bash
./vad_tester -threshold 0.3 -neg-threshold 0.1 -min-speech 150
```

2. Conservative detection with longer segments (reduces false positives):
```bash
./vad_tester -threshold 0.7 -min-silence 1000 -min-speech 500
```

3. Fine-tuned for quick responses (good for short commands):
```bash
./vad_tester -threshold 0.5 -neg-threshold 0.4 -min-silence 200 -min-speech 200 -speech-pad 10
```

4. Optimized for noisy environments:
```bash
./vad_tester -threshold 0.6 -neg-threshold 0.45 -min-silence 700 -min-speech 300 -speech-pad 50
```

## Testing and Comparison

### Using Compare Scripts

The repository includes two comparison scripts that help you evaluate different parameter combinations:

1. Using `compare_params.sh`:
```bash
# Make the script executable
chmod +x compare_params.sh

# Run all tests
./compare_params.sh

# View results in output.txt
cat output.txt
```

2. Using `compare_params2.sh` (for second sample):
```bash
chmod +x compare_params2.sh
./compare_params2.sh
```

### Test Cases Included

Each comparison script runs 11 different test cases:
1. Base parameters (default settings)
2. Low speech threshold (more sensitive)
3. High speech threshold (more conservative)
4. Low negative threshold (slower silence detection)
5. High negative threshold (faster silence detection)
6. Short minimum silence (more segments)
7. Long minimum silence (fewer segments)
8. Short minimum speech (catches brief sounds)
9. Long minimum speech (filters short sounds)
10. No padding (precise boundaries)
11. Large padding (includes context)

### Analyzing Test Results

The test output includes:
```
Test: [Test Name]
------------------------------------------------------------
Parameters:
  Threshold: [value]
  Negative Threshold: [value]
  Min Silence Duration: [value] ms
  Min Speech Duration: [value] ms
  Speech Pad: [value] ms
------------------------------------------------------------
Detected speech segments:
1. [start] - [end] ([duration] sec)
...
Total audio duration: [total] sec
Total speech duration: [speech] sec ([percentage]%)
```

### Choosing the Best Parameters

1. For general use:
   - Start with default parameters
   - Run both comparison scripts
   - Look for the test case that gives the best balance of accuracy and segment count

2. For specific use cases:
   - Quiet speech: Use Test 2 (Low speech threshold)
   - Noisy environment: Use Test 3 (High speech threshold)
   - Quick commands: Use Test 6 (Short minimum silence)
   - Long phrases: Use Test 7 (Long minimum silence)

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
Real-time factor: 0.15
```

## Troubleshooting

Common issues and solutions:

1. No speech detected:
   - Try lowering the speech threshold
   - Check if audio format is correct (PCM)
   - Verify sample rate matches audio

2. Too many segments:
   - Increase minimum silence duration
   - Increase minimum speech duration
   - Raise speech threshold

3. Segments cut too early/late:
   - Adjust speech padding
   - Fine-tune negative threshold
   - Check for background noise

4. Poor performance:
   - Verify ONNX Runtime installation
   - Check CPU usage
   - Monitor memory consumption 