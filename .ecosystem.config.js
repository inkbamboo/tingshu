const os = process.env.OS ? process.env.OS : "darwin";
const watch = process.env.WATCH ? process.env.WATCH : false;
const env = process.env.ENV ? process.env.ENV : 'dev';
const go_env = process.env.ENV ? process.env.ENV : 'dev';
const configs = {
    instances: 1,
    max_memory_restart: "1000M",
    autorestart: false,
    max_restarts: 10,
    exec_mode: "fork",
    interpreter: "none",
    log_file: '../../logs/tingshu/tingshu.log',
    combine_logs: true,
    watch,
    env_uat: {
        env: 'local',
        GO_ENV: go_env
    },
    env_stage: {
        env: 'stage',
        GO_ENV: go_env
    },
    env_stress: {
        env: 'stress',
        GO_ENV: go_env
    },
    env_prod: {
        env: 'prod',
        GO_ENV: go_env
    }
}
module.exports = {
    apps: [
        {
            name: "tingshu-server",
            cwd: `./build/tingshu-server_${os}_amd64`,
            script: "tingshu-server",
            args:`--env ${env} s`,
            ...configs,
        }
    ],
};