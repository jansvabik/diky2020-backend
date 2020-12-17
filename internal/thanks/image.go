package thanks

import (
	"fmt"
	"image/jpeg"
	"os"

	"github.com/beevik/guid"
	"github.com/gofiber/fiber/v2"
	"github.com/nfnt/resize"
	"github.com/noltio/diky2020-backend/pkg/server"
)

// ImageUploadHandler handles requests for thanks image upload
func ImageUploadHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if file == nil {
		c.Locals("uploadedFile", false)
		return c.Next()
	}
	if err != nil {
		fmt.Println(err.Error())
		return server.APIInternalServerError(c)
	}

	// filesize condition
	if file.Size > 8388608 {
		return server.APIError(c, "Maximální velikost nahrávaného obrázku je 8 MB.", 400)
	}

	// mime type condition
	ct := file.Header["Content-Type"][0]
	mimes := map[string]string{
		"image/jpeg": "jpg",
		"image/png":  "png",
		"image/gif":  "gif",
	}
	if _, ok := mimes[ct]; !ok {
		return server.APIError(c, "Povolené formáty obrázků jsou pouze .jpg, .jpeg, .png a .gif.", 400)
	}

	// filename (for both temp and resized file)
	filename := fmt.Sprintf("%s.%s", guid.New().String(), mimes[ct])

	// save the original file
	c.SaveFile(file, filename)

	// open the uploaded file
	nrf, err := os.Open(filename)
	if err != nil {
		fmt.Println(err.Error())
		return server.APIInternalServerError(c)
	}
	defer nrf.Close()

	// decode jpeg into image.Image
	img, err := jpeg.Decode(nrf)
	if err != nil {
		fmt.Println(err.Error())
		return server.APIInternalServerError(c)
	}

	// resize to width 600 using Lanczos resampling
	// and preserve the original aspect ratio
	m := resize.Resize(600, 0, img, resize.Lanczos3)

	// save new file
	out, err := os.Create(filename)
	if err != nil {
		fmt.Println(err.Error())
		return server.APIInternalServerError(c)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)

	// store the data in local variables and move to the next handler
	c.Locals("uploadedFile", true)
	c.Locals("uploadedFileName", filename)
	return c.Next()
}
