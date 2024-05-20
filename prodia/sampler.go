package prodia

const (
	SamplerDPMPPSDEKarras       = "DPM++ SDE Karras"
	SamplerDPMPPSDEExponential  = "DPM++ SDE Exponential"
	SamplerEuler                = "Euler"
	SamplerEulerA               = "Euler a"
	SamplerLMS                  = "LMS"
	SamplerHeun                 = "Heun"
	SamplerDPM2                 = "DPM2"
	SamplerDPM2a                = "DPM2 a"
	SamplerDPMPP2MSDEHeunKarras = "DPM++ 2M SDE Heun Karras"
	SamplerLMSKarras            = "LMS Karras"
	SamplerDPM2Karras           = "DPM2 Karras"
	SamplerDDIM                 = "DDIM"
	SamplerUniPC                = "UniPC"
)

var mapStyleSampler = map[string]string{
	"realistic":     SamplerDPMPPSDEKarras,
	"cartoon":       SamplerEulerA,
	"chibi":         SamplerEulerA,
	"abstract":      SamplerDPMPPSDEExponential,
	"minimalist":    SamplerHeun,
	"vintage":       SamplerLMS,
	"fantasy":       SamplerDPMPPSDEKarras,
	"surreal":       SamplerDPMPPSDEExponential,
	"pop-art":       SamplerDPMPP2MSDEHeunKarras,
	"watercolor":    SamplerDDIM,
	"pixel-art":     SamplerEuler,
	"line-art":      SamplerHeun,
	"cyberpunk":     SamplerDPMPPSDEKarras,
	"steampunk":     SamplerDPMPPSDEKarras,
	"art-deco":      SamplerLMS,
	"gothic":        SamplerDPM2Karras,
	"impressionist": SamplerDPM2,
	"expressionist": SamplerDPM2a,
	"sci-fi":        SamplerDPMPP2MSDEHeunKarras,
	"3d-render":     SamplerUniPC,
	"retro":         SamplerLMSKarras,
}

func GetSampler(style string) string {
	sampler, exists := mapStyleSampler[style]
	if !exists {
		return SamplerEulerA
	}
	return sampler
}
