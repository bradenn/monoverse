#version 150

in vec3 vertNorm;

uniform vec3 lightDir;
uniform vec3 lightCol;

out vec4 outColor;

void main() {
    float brightness = max(dot(vertNorm, -normalize(lightDir)), 0.0);

    outColor = vec4(0.0) + brightness * vec4(lightCol, 1.0);
}