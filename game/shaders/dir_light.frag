layout(location = 0) out vec4 fbColor0Out;

uniform sampler2D fbColor0TextureIn;
uniform sampler2D fbColor1TextureIn;
uniform sampler2D fbDepthTextureIn;
uniform sampler2DShadow fbShadowTextureIn;

layout (std140) uniform Camera
{
	mat4 projectionMatrixIn;
	mat4 viewMatrixIn;
	mat4 cameraMatrixIn;
};

layout (std140) uniform Light
{
	mat4 lightProjectionMatrixIn;
	mat4 lightViewMatrixIn;
	mat4 lightMatrixIn;
};

uniform vec3 lightIntensityIn;

smooth in vec2 texCoordInOut;

const float pi = 3.141592;

struct fresnelInput {
	vec3 reflectanceF0;
	vec3 halfDirection;
	vec3 lightDirection;
};

vec3 calculateFresnel(fresnelInput i) {
	float halfLightDot = clamp(dot(i.halfDirection, i.lightDirection), 0.0, 1.0);
	return i.reflectanceF0 + (1.0 - i.reflectanceF0) * pow(1.0 - halfLightDot, 5.0);
}

struct distributionInput {
	float roughness;
	vec3 normal;
	vec3 halfDirection;
};

float calculateDistribution(distributionInput i) {
	float sqrRough = i.roughness * i.roughness;
	float halfNormDot = dot(i.normal, i.halfDirection);
	float denom = halfNormDot * halfNormDot * (sqrRough - 1.0) + 1.0;
	return sqrRough / (pi * denom * denom);
}

struct geometryInput {
	float roughness;
};

float calculateGeometry(geometryInput i) {
	// TODO: Use better model
	return 1.0 / 4.0;
}

struct directionalSetup {
	float roughness;
	vec3 reflectedColor;
	vec3 refractedColor;
	vec3 viewDirection;
	vec3 lightDirection;
	vec3 normal;
	vec3 lightIntensity;
};

vec3 calculateDirectionalHDR(directionalSetup s) {
	vec3 halfDirection = normalize(s.lightDirection + s.viewDirection);
	vec3 fresnel = calculateFresnel(fresnelInput(
		s.reflectedColor,
		halfDirection,
		s.lightDirection
	));
	float distributionFactor = calculateDistribution(distributionInput(
		s.roughness,
		s.normal,
		halfDirection
	));
	float geometryFactor = calculateGeometry(geometryInput(
		s.roughness
	));
	vec3 reflectedHDR = fresnel * distributionFactor * geometryFactor;
	vec3 refractedHDR = (vec3(1.0) - fresnel) * s.refractedColor / pi;
	return (reflectedHDR + refractedHDR) * s.lightIntensity * clamp(dot(s.normal, s.lightDirection), 0.0, 1.0);
}

float textureClampToBorder(sampler2DShadow tex, vec3 coord, float dValue) {
	if (coord.x < 0.0 || coord.x > 1.0 || coord.y < 0.0 || coord.y > 1.0) {
		return dValue;
	}
	return texture(tex, coord);
}

void main()
{
	vec3 ndcPosition = vec3(
		(texCoordInOut.x - 0.5) * 2.0,
		(texCoordInOut.y - 0.5) * 2.0,
		texture(fbDepthTextureIn, texCoordInOut).x * 2.0 - 1.0
	);
	vec3 clipPosition = vec3(
		ndcPosition.x / projectionMatrixIn[0][0],
		ndcPosition.y / projectionMatrixIn[1][1],
		-1.0
	);
	vec3 viewPosition = clipPosition * projectionMatrixIn[3][2] / (projectionMatrixIn[2][2] + ndcPosition.z);
	vec3 worldPosition = (cameraMatrixIn * vec4(viewPosition, 1.0)).xyz;
	vec3 cameraPosition = cameraMatrixIn[3].xyz;

	vec4 albedoMetalness = texture(fbColor0TextureIn, texCoordInOut);
	vec4 normalRoughness = texture(fbColor1TextureIn, texCoordInOut);
	vec3 baseColor = albedoMetalness.xyz;
	vec3 normal = normalize(normalRoughness.xyz);
	float metalness = albedoMetalness.w;
	float roughness = normalRoughness.w;

	vec3 refractedColor = baseColor * (1.0 - metalness);
	vec3 reflectedColor = mix(vec3(0.02), baseColor, metalness);

	vec3 lightDirection = normalize(lightMatrixIn[2].xyz);

	vec3 hdr = calculateDirectionalHDR(directionalSetup(
		roughness,
		reflectedColor,
		refractedColor,
		normalize(cameraPosition - worldPosition),
		lightDirection,
		normal,
		lightIntensityIn
	));

	vec4 lightPosition = lightProjectionMatrixIn * lightViewMatrixIn * vec4(worldPosition, 1.0);
	float directness = clamp(dot(normal, lightDirection), 0.0, 1.0);
	lightPosition.xyz = lightPosition.xyz * 0.5 + 0.5;
	lightPosition.z /= lightPosition.w;
	lightPosition.z -= 0.0005;

	vec2 shift = 1.0 / vec2(textureSize(fbShadowTextureIn, 0));

	vec3 shifts[9];
	shifts[0] =	vec3(0.0, 0.0, 0.0);
	shifts[1] =	vec3(-1.0, 0.0, 0.0);
	shifts[2] =	vec3(1.0, 0.0, 0.0);
	shifts[3] =	vec3(0.0, -1.0, 0.0);
	shifts[4] =	vec3(0.0, 1.0, 0.0);
	shifts[5] =	vec3(-1.0, -1.0, 0.0);
	shifts[6] =	vec3(1.0, -1.0, 0.0);
	shifts[7] =	vec3(-1.0, 1.0, 0.0);
	shifts[8] =	vec3(1.0, 1.0, 0.0);

	float amount = 0.0;
	for (int i = 0; i < 9; i++) {
		highp vec3 texPos = lightPosition.xyz + shifts[i] * vec3(shift.x, shift.y, 1.0);
		float probability = textureClampToBorder(fbShadowTextureIn, texPos, 1.0);
		amount = max(amount, probability);
	}

	float factor = (clamp(directness, 0.3, 0.5) - 0.3) / 0.2;
	amount = amount * factor;

	fbColor0Out = vec4(hdr * amount, 1.0);
}
