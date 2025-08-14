# ffmpeg 查看文件信息
ffmpeg -i audio.wav

# ffmpeg 采用率转换
ffmpeg -i audio.wav -ar 16000 audio_16k.wav

# ffmpeg 播放 16k pcm音频
ffplay -f s16le -ar 16000 -ac 1 your_audio.pcm

# ffplay 播放 16k pcm音频 (使用这个)
ffplay -ar 16000 -acodec pcm_s16le -f s16le -i your_audio.pcm