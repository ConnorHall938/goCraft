#version 460 core

layout(location = 0) in vec3 aPos;   // x, y, z
layout(location = 1) in vec2 aUV;    // u, v
layout(location = 2) in vec3 aTint;  // r, g, b

out vec2 vUV;
out vec3 vTint;

uniform mat4 uMVP;

void main() {
    vUV = aUV;
    vTint = aTint;
    gl_Position = uMVP * vec4(aPos, 1.0);
}
