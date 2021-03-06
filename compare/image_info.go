package compare

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/corona10/goimagehash"
	"gopkg.in/gographics/imagick.v2/imagick"
)

const (
	IMAGE_ERROR     = "image_error"
	WIDTH           = "width"
	HEIGHT          = "height"
	QUALITY         = "quality"
	FRAME_NUM       = "frame_num"
	AVERAGE_HASH    = "average_hash"
	DIFFERENCE_HASH = "difference_hash"
	PERCEPTION_HASH = "perception_hash"
)

type ImageInfo struct {
	Error string

	Width    int
	Height   int
	Quality  int
	FrameNum int

	AverageHash    *goimagehash.ImageHash
	DifferenceHash *goimagehash.ImageHash
	PerceptionHash *goimagehash.ImageHash
}

func isImage(format string) bool {
	// FIXME 检查格式判断是否为图片
	return true
}

func NewImageInfo() *ImageInfo {
	return &ImageInfo{
		AverageHash:    goimagehash.NewImageHash(100, goimagehash.Unknown),
		DifferenceHash: goimagehash.NewImageHash(100, goimagehash.Unknown),
		PerceptionHash: goimagehash.NewImageHash(100, goimagehash.Unknown),
	}
}

func ParseImageInfo(data []byte) *ImageInfo {
	info := NewImageInfo()
	mw := imagick.NewMagickWand()
	if err := mw.ReadImageBlob(data); err != nil {
		info.Error = err.Error()
		return info
	}
	info.Width = int(mw.GetImageWidth())
	info.Height = int(mw.GetImageHeight())
	info.Quality = int(mw.GetCompressionQuality())
	if mw.GetImageFormat() == "GIF" {
		frameNum, _ := mw.GetImageLength()
		info.FrameNum = int(frameNum)
	}
	r := bytes.NewReader(data)
	img, _, err := image.Decode(r)
	if err != nil {
		info.Error = err.Error()
		return info
	}
	info.AverageHash, _ = goimagehash.AverageHash(img)
	info.DifferenceHash, _ = goimagehash.DifferenceHash(img)
	info.PerceptionHash, _ = goimagehash.PerceptionHash(img)
	return info
}

func (this *ImageInfo) Compare(other *ImageInfo, result map[string]Diff) {
	if this.Error != other.Error {
		result[IMAGE_ERROR] = &StringDiff{this.Error, other.Error}
	}
	if this.Width != other.Width {
		result[WIDTH] = &IntDiff{this.Width, other.Width}
	}
	if this.Height != other.Height {
		result[HEIGHT] = &IntDiff{this.Height, other.Height}
	}
	if this.Quality != other.Quality {
		result[QUALITY] = &IntDiff{this.Quality, other.Quality}
	}
	if this.FrameNum != other.FrameNum {
		result[FRAME_NUM] = &IntDiff{this.FrameNum, other.FrameNum}
	}
	if dis, _ := this.AverageHash.Distance(other.AverageHash); dis != 0 {
		result[AVERAGE_HASH] = &IntDiff{dis, 0}
	}
	if dis, _ := this.DifferenceHash.Distance(other.DifferenceHash); dis != 0 {
		result[DIFFERENCE_HASH] = &IntDiff{dis, 0}
	}
	if dis, _ := this.PerceptionHash.Distance(other.PerceptionHash); dis != 0 {
		result[PERCEPTION_HASH] = &IntDiff{dis, 0}
	}
}
