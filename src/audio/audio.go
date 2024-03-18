package audio

import (
	"bytes"
	"kaiju/audio/audio_system"
	"log/slog"
	"time"

	"github.com/ebitengine/oto/v3"
)

type Audio struct {
	otoCtx  *oto.Context
	options oto.NewContextOptions
}

func NewAudio() (Audio, error) {
	a := Audio{
		options: oto.NewContextOptions{},
	}
	a.options.SampleRate = 48000
	a.options.ChannelCount = 2
	a.options.Format = oto.FormatFloat32LE
	otoCtx, readyChan, err := oto.NewContext(&a.options)
	if err != nil {
		return Audio{}, err
	}
	a.otoCtx = otoCtx
	<-readyChan
	return a, nil
}

func (a *Audio) Play(wav *audio_system.Wav) {
	if wav == nil {
		slog.Error("Wav is nil")
		return
	}
	data := wav.WavData
	// TODO:  Rather than doing this real-time, it should be a part of the
	// import process of the asset.  This is a temporary solution.
	if int(wav.Channels) != a.options.ChannelCount {
		slog.Warn("Rechanneling audio, this is a temporary solution",
			slog.Int("channels", int(wav.Channels)),
			slog.Int("target", int(a.options.ChannelCount)))
		data = audio_system.Rechannel(wav, int16(a.options.ChannelCount))
	}
	if int(wav.SampleRate) != a.options.SampleRate {
		slog.Warn("Resampling audio, this is a temporary solution",
			slog.Int("sampleRate", int(wav.SampleRate)),
			slog.Int("target", int(a.options.SampleRate)))
		data = audio_system.Resample(wav, int32(a.options.SampleRate))
	}
	if wav.FormatType != audio_system.WavFormatFloat {
		slog.Warn("Converting audio to float, this is a temporary solution",
			slog.String("format", "Float"))
		data = audio_system.Pcm2Float(wav)
	}
	player := a.otoCtx.NewPlayer(bytes.NewReader(data))
	player.Play()
	go func() {
		for player.IsPlaying() {
			time.Sleep(time.Millisecond)
		}
		if err := player.Close(); err != nil {
			slog.Error(err.Error())
		}
	}()
}
