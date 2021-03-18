package compare

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/corona10/goimagehash"
	"github.com/gographics/imagick/imagick"
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
	info.ParseImageInfo(data)
	return info

	if !ParseImageBaseInfo(info, data) || !ParseImageHash(info, data) {
		return info
	}
	return info
}

func (this *ImageInfo) ParseImageInfo(data []byte) {
	mw := imagick.NewMagickWand()
	if err := mw.ReadImageBlob(data); err != nil {
		info.Error = err.Error()
		return
	}
	this.Width = mw.GetImageWidth()
	this.Height = mw.GetImageHeight()
	this.Quality = mw.GetCompressionQuality()
	if mw.GetImageFormat() == "GIF" {
		this.FrameNum, _ = mw.GetImageLength()
	}
	r := bytes.NewReader(data)
	img, _, err := image.Decode(r)
	if err != nil {
		info.Error = err.Error()
		return
	}
	this.AverageHash, _ = goimagehash.AverageHash(img)
	this.DifferenceHash, _ = goimagehash.DifferenceHash(img)
	this.PerceptionHash, _ = goimagehash.PerceptionHash(img)
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
