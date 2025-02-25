package speech

import (
	"encoding/binary"
	"log/slog"
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDetectorConfigIsValid(t *testing.T) {
	tcs := []struct {
		name string
		cfg  DetectorConfig
		err  string
	}{
		{
			name: "missing ModelPath",
			cfg: DetectorConfig{
				ModelPath: "",
			},
			err: "invalid ModelPath: should not be empty",
		},
		{
			name: "invalid SampleRate",
			cfg: DetectorConfig{
				ModelPath:  "../testfiles/silero_vad.onnx",
				SampleRate: 48000,
			},
			err: "invalid SampleRate: valid values are 8000 and 16000",
		},
		{
			name: "invalid Threshold",
			cfg: DetectorConfig{
				ModelPath:  "../testfiles/silero_vad.onnx",
				SampleRate: 16000,
				Threshold:  0,
			},
			err: "invalid Threshold: should be in range (0, 1)",
		},
		{
			name: "invalid MinSilenceDurationMs",
			cfg: DetectorConfig{
				ModelPath:            "../testfiles/silero_vad.onnx",
				SampleRate:           16000,
				Threshold:            0.5,
				MinSilenceDurationMs: -1,
			},
			err: "invalid MinSilenceDurationMs: should be a positive number",
		},
		{
			name: "invalid SpeechPadMs",
			cfg: DetectorConfig{
				ModelPath:   "../testfiles/silero_vad.onnx",
				SampleRate:  16000,
				Threshold:   0.5,
				SpeechPadMs: -1,
			},
			err: "invalid SpeechPadMs: should be a positive number",
		},
		{
			name: "invalid MinSpeechDurationMs",
			cfg: DetectorConfig{
				ModelPath:          "../testfiles/silero_vad.onnx",
				SampleRate:         16000,
				Threshold:          0.5,
				MinSpeechDurationMs: -1,
			},
			err: "invalid MinSpeechDurationMs: should be a positive number",
		},
		{
			name: "invalid NegativeThreshold range",
			cfg: DetectorConfig{
				ModelPath:         "../testfiles/silero_vad.onnx",
				SampleRate:        16000,
				Threshold:         0.5,
				NegativeThreshold: 0,
			},
			err: "invalid NegativeThreshold: should be in range (0, 1)",
		},
		{
			name: "invalid NegativeThreshold greater than Threshold",
			cfg: DetectorConfig{
				ModelPath:         "../testfiles/silero_vad.onnx",
				SampleRate:        16000,
				Threshold:         0.5,
				NegativeThreshold: 0.6,
			},
			err: "invalid NegativeThreshold: should be less than Threshold",
		},
		{
			name: "valid",
			cfg: DetectorConfig{
				ModelPath:         "../testfiles/silero_vad.onnx",
				SampleRate:        16000,
				Threshold:         0.5,
				NegativeThreshold: 0.35,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.cfg.IsValid()
			if tc.err != "" {
				require.EqualError(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewDetector(t *testing.T) {
	cfg := DetectorConfig{
		ModelPath:  "../testfiles/silero_vad.onnx",
		SampleRate: 16000,
		Threshold:  0.5,
	}

	sd, err := NewDetector(cfg)
	require.NoError(t, err)
	require.NotNil(t, sd)

	err = sd.Destroy()
	require.NoError(t, err)
}

func TestSpeechDetection(t *testing.T) {
	cfg := DetectorConfig{
		ModelPath:  "../testfiles/silero_vad.onnx",
		SampleRate: 16000,
		Threshold:  0.5,
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	sd, err := NewDetector(cfg)
	require.NoError(t, err)
	require.NotNil(t, sd)
	defer func() {
		require.NoError(t, sd.Destroy())
	}()

	readSamplesFromFile := func(path string) []float32 {
		data, err := os.ReadFile(path)
		require.NoError(t, err)

		samples := make([]float32, 0, len(data)/4)
		for i := 0; i < len(data); i += 4 {
			samples = append(samples, math.Float32frombits(binary.LittleEndian.Uint32(data[i:i+4])))
		}
		return samples
	}

	samples := readSamplesFromFile("../testfiles/samples.pcm")
	samples2 := readSamplesFromFile("../testfiles/samples2.pcm")

	t.Run("detect", func(t *testing.T) {
		segments, err := sd.Detect(samples)
		require.NoError(t, err)
		require.NotEmpty(t, segments)
		require.Equal(t, []Segment{
			{
				SpeechStartAt: 1.056,
				SpeechEndAt:   1.632,
			},
			{
				SpeechStartAt: 2.88,
				SpeechEndAt:   3.232,
			},
			{
				SpeechStartAt: 4.448,
				SpeechEndAt:   0,
			},
		}, segments)

		err = sd.Reset()
		require.NoError(t, err)

		segments, err = sd.Detect(samples2)
		require.NoError(t, err)
		require.NotEmpty(t, segments)
		require.Equal(t, []Segment{
			{
				SpeechStartAt: 3.008,
				SpeechEndAt:   6.24,
			},
			{
				SpeechStartAt: 7.072,
				SpeechEndAt:   8.16,
			},
		}, segments)
	})

	t.Run("reset", func(t *testing.T) {
		err = sd.Reset()
		require.NoError(t, err)

		segments, err := sd.Detect(samples)
		require.NoError(t, err)
		require.NotEmpty(t, segments)
		require.Equal(t, []Segment{
			{
				SpeechStartAt: 1.056,
				SpeechEndAt:   1.632,
			},
			{
				SpeechStartAt: 2.88,
				SpeechEndAt:   3.232,
			},
			{
				SpeechStartAt: 4.448,
				SpeechEndAt:   0,
			},
		}, segments)
	})

	t.Run("speech padding", func(t *testing.T) {
		cfg.SpeechPadMs = 10
		sd, err := NewDetector(cfg)
		require.NoError(t, err)
		require.NotNil(t, sd)
		defer func() {
			require.NoError(t, sd.Destroy())
		}()

		segments, err := sd.Detect(samples)
		require.NoError(t, err)
		require.NotEmpty(t, segments)
		require.Equal(t, []Segment{
			{
				SpeechStartAt: 1.056 - 0.01,
				SpeechEndAt:   1.632 + 0.01,
			},
			{
				SpeechStartAt: 2.88 - 0.01,
				SpeechEndAt:   3.232 + 0.01,
			},
			{
				SpeechStartAt: 4.448 - 0.01,
				SpeechEndAt:   0,
			},
		}, segments)
	})

	t.Run("negative threshold", func(t *testing.T) {
		cfg.SpeechPadMs = 0
		cfg.NegativeThreshold = 0.3 // Lower negative threshold should result in longer speech segments
		sd, err := NewDetector(cfg)
		require.NoError(t, err)
		require.NotNil(t, sd)
		defer func() {
			require.NoError(t, sd.Destroy())
		}()

		segments, err := sd.Detect(samples)
		require.NoError(t, err)
		require.NotEmpty(t, segments)

		// Test that we can change the negative threshold after creation
		sd.SetNegativeThreshold(0.2)
		err = sd.Reset()
		require.NoError(t, err)
		
		segments2, err := sd.Detect(samples)
		require.NoError(t, err)
		require.NotEmpty(t, segments2)
		
		// With an even lower threshold, we expect longer speech segments
		// or potentially fewer segments due to merging
		require.True(t, len(segments2) <= len(segments), 
			"Expected fewer or equal number of segments with lower threshold")
	})

	t.Run("min speech duration", func(t *testing.T) {
		// Reset config
		cfg.SpeechPadMs = 0
		cfg.NegativeThreshold = 0
		
		// First run with no minimum speech duration
		cfg.MinSpeechDurationMs = 0
		sd, err := NewDetector(cfg)
		require.NoError(t, err)
		require.NotNil(t, sd)
		defer func() {
			require.NoError(t, sd.Destroy())
		}()

		segments, err := sd.Detect(samples)
		require.NoError(t, err)
		require.NotEmpty(t, segments)
		initialSegmentCount := len(segments)

		// Now set a high minimum speech duration that should filter out some segments
		sd.SetMinSpeechDurationMs(1000) // 1 second
		err = sd.Reset()
		require.NoError(t, err)
		
		segments2, err := sd.Detect(samples)
		require.NoError(t, err)
		
		// With a higher minimum speech duration, we expect fewer segments
		require.True(t, len(segments2) <= initialSegmentCount, 
			"Expected fewer segments with higher minimum speech duration")
	})
}
