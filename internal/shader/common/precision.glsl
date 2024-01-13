#if defined(GL_FRAGMENT_PRECISION_HIGH)
precision highp float;
precision highp sampler2DShadow;
#else
precision mediump float;
precision mediump sampler2DShadow;
#endif