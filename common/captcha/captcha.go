package captcha

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"strconv"
	"time"
)

type Captcha struct {
	config *Config
}

type Config struct {
	Width      int
	Height     int
	CharsCount int
	FontSize   int
}

func New(config *Config) *Captcha {
	return &Captcha{config: config}
}

// 生成验证码
func (c *Captcha) Generate() ([]byte, string) {

	charset := "0123456789"
	rand.Seed(time.Now().UnixNano())

	chars := make([]rune, c.config.CharsCount)
	for i := range chars {
		chars[i] = rune(charset[rand.Intn(len(charset))])
	}

	var svg bytes.Buffer
	svg.WriteString("<svg xmlns='http://www.w3.org/2000/svg' width='" + strconv.Itoa(c.config.Width) + "' height='" + strconv.Itoa(c.config.Height) + "'>")
	svg.WriteString("<rect width='" + strconv.Itoa(c.config.Width) + "' height='" + strconv.Itoa(c.config.Height) + "' fill='white'/>")

	// 从左下角开始绘制文本
	// 假设文本的字体大小为 c.config.FontSize
	// 将文本的起始点设置为 (0, c.config.Height - c.config.FontSize)
	// 为了保持文本居中，我们需要添加字符间距，这里假设每个字符之间有字符宽度的间距
	charWidth := c.config.FontSize / 2
	for i, char := range chars {
		// 计算文本的 x 和 y 坐标
		x := float64(i)*float64(charWidth) + float64(c.config.FontSize)
		y := float64(c.config.Height - c.config.FontSize) + float64(c.config.FontSize) / 2
		// 获取随机颜色并格式化为十六进制字符串
		col := randomColor().(color.RGBA)
		r, g, b := col.R, col.G, col.B
		svgColor := fmt.Sprintf("#%02x%02x%02x", r, g, b)
		svg.WriteString(fmt.Sprintf("<text x='%f' y='%f' font-family='Verdana' font-size='%d' fill='%s'>%c</text>",
			x, y, c.config.FontSize, svgColor, char))
	}

	// 生成干扰线
	svg.WriteString("<g fill='none' stroke='black'>")
	for i := 0; i < 4; i++ {
		startX := rand.Intn(c.config.Width)
		startY := rand.Intn(c.config.Height)
		endX := rand.Intn(c.config.Width)
		endY := rand.Intn(c.config.Height)
		svg.WriteString(fmt.Sprintf("<line x1='%d' y1='%d' x2='%d' y2='%d' stroke='%02x%02x%02x'/>", startX, startY, endX, endY, 255-int(randomColor().(color.RGBA).R), 255-int(randomColor().(color.RGBA).G), 255-int(randomColor().(color.RGBA).B)))
	}
	svg.WriteString("</g>")
	svg.WriteString("</svg>")

	return svg.Bytes(), string(chars)
}

// 在图像上绘制随机线条作为噪声
func (c *Captcha) DrawNoise(img *image.RGBA) {
	for i := 0; i < 3; i++ {
		startX := rand.Intn(c.config.Width)
		startY := rand.Intn(c.config.Height)
		endX := rand.Intn(c.config.Width)
		endY := rand.Intn(c.config.Height)
		drawLine(img, image.Point{startX, startY}, image.Point{endX, endY}, randomColor())
	}
}

// 生成一个随机的RGBA颜色
func randomColor() color.Color {
	r := rand.Intn(256)
	g := rand.Intn(256)
	b := rand.Intn(256)
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(255)}
}

// 绘制干扰线
func drawLine(img *image.RGBA, p1, p2 image.Point, c color.Color) {
	var dx, dy int
	if p1.X < p2.X {
		dx = p2.X - p1.X
	} else {
		dx = p1.X - p2.X
	}
	if p1.Y < p2.Y {
		dy = p2.Y - p1.Y
	} else {
		dy = p1.Y - p2.Y
	}
	if dx > dy {
		for x := 0; x <= dx; x++ {
			y := int(float64(dy) * float64(x) / float64(dx))
			img.Set(p1.X+x, p1.Y+y, c)
		}
	} else {
		for y := 0; y <= dy; y++ {
			x := int(float64(dx) * float64(y) / float64(dy))
			img.Set(p1.X+x, p1.Y+y, c)
		}
	}
}
