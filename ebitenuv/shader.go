package ebitenuv

const shaderSource = `
//kage:unit pixels

package main

func Fragment(dstPos vec4, srcPos vec2, color vec4) vec4 {
	srcColor := imageSrc0At(srcPos)
	if srcColor.a == 0 {
		return vec4(0)
	}

	// R,G encode UV pixel coordinates (0-255). Kage normalizes to 0.0-1.0.
	uvX := srcColor.r * 255.0
	uvY := srcColor.g * 255.0

	// Sample center of target pixel in lookup
	lookupColor := imageSrc1At(vec2(uvX + 0.5, uvY + 0.5))
	if lookupColor.a == 0 {
		return vec4(0)
	}

	return lookupColor
}
`
