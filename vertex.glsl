#version 150

in vec3 position;
in vec3 normal;

uniform mat4 model;
uniform mat4 view;
uniform mat4 proj;

out vec3 vertNorm;

void main() {
    gl_Position = proj * view * model * vec4(position, 1.0);
    vertNorm = (model * vec4(normal, 1.0)).xyz;
}