/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  swcMinify: true,
  experimental: {
    appDir: true
  },
  output: 'standalone',
  distDir: '.next',
  generateEtags: true,
  poweredByHeader: false,
  compress: true,
  webpack: (config, { isServer }) => {
    if (!isServer) {
      config.resolve.fallback = {
        ...config.resolve.fallback,
        fs: false,
        net: false,
        tls: false,
        crypto: false,
        stream: false,
        path: false,
        zlib: false,
        http: false,
        https: false,
        os: false,
        url: false,
        assert: false,
        buffer: false,
        process: false,
        util: false,
        events: false,
        punycode: false,
        querystring: false,
        string_decoder: false,
        sys: false,
        timers: false,
        dns: false,
        dgram: false,
        child_process: false,
        cluster: false,
        module: false,
        readline: false,
        repl: false,
        tty: false,
        vm: false,
        constants: false,
        domain: false,
        inspector: false,
        perf_hooks: false,
        v8: false,
        worker_threads: false,
      };
    }
    return config;
  },
  env: {
    PORT: process.env.PORT || '3000'
  },
  serverRuntimeConfig: {
    port: process.env.PORT || '3000'
  },
  publicRuntimeConfig: {
    port: process.env.PORT || '3000'
  }
}

module.exports = nextConfig; 