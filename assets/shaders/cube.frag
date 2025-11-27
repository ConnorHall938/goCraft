#version 460 core

in vec2 vUV;
in vec3 vTint;

uniform sampler2D uAtlas;  // atlas texture

out vec4 FragColor;

void main() {
    vec4 texColor = texture(uAtlas, vUV);
    FragColor = vec4(texColor.rgb * vTint, texColor.a);
}
