package prodia

import "fmt"

const (
	ModelAnimagineXLV3 = "animagineXLV3_v30.safetensors [75f2f05b]"
	ModelDreamshaperXL = "dreamshaperXL10_alpha2.safetensors [c8afe2ef]"
	ModelDynavisionXL  = "dynavisionXL_0411.safetensors [c39cc051]"
	ModelRealismEngine = "realismEngineSDXL_v10.safetensors [af771c3f]"
	ModelRealvisXL     = "realvisxlV40.safetensors [f7fdcb51]"
)

var mapStyleModel = map[string]string{
	"realistic":     ModelRealismEngine,
	"cartoon":       ModelAnimagineXLV3,
	"chibi":         ModelAnimagineXLV3,
	"abstract":      ModelDreamshaperXL,
	"minimalist":    ModelDreamshaperXL,
	"vintage":       ModelRealvisXL,
	"fantasy":       ModelAnimagineXLV3,
	"surreal":       ModelDreamshaperXL,
	"pop-art":       ModelAnimagineXLV3,
	"watercolor":    ModelDreamshaperXL,
	"pixel-art":     ModelAnimagineXLV3,
	"line-art":      ModelDreamshaperXL,
	"cyberpunk":     ModelAnimagineXLV3,
	"steampunk":     ModelAnimagineXLV3,
	"art-deco":      ModelDreamshaperXL,
	"gothic":        ModelDreamshaperXL,
	"impressionist": ModelDreamshaperXL,
	"expressionist": ModelDreamshaperXL,
	"sci-fi":        ModelAnimagineXLV3,
	"3d-render":     ModelDynavisionXL,
	"retro":         ModelRealvisXL,
}

func GetModel(style string) string {
	model, exists := mapStyleModel[style]
	if !exists {
		return ModelDreamshaperXL
	}
	return model
}

var mapProductSize = map[string]string{
	"t-shirt":    "819x1024",
	"tumbler":    "1024x1024",
	"phone-case": "576x1024",
	"hoodie":     "819x1024",
	"mug":        "1024x1024",
	"tote-bag":   "819x1024",
	"pillow":     "1024x1024",
	"poster":     "683x1024",
	"notebook":   "1024x683",
	"sticker":    "1024x1024",
}

func GetSize(product string) (int, int) {
	size, exists := mapProductSize[product]
	if !exists {
		return 1024, 1024
	}
	var width, height int
	_, _ = fmt.Sscanf(size, "%dx%d", &width, &height)
	return width, height
}
