INSERT INTO projects.frameworks (name, base_image, root_dir, output_dir, install_cmd, build_cmd, run_cmd)
VALUES
    (
        'Node.js',
        'node:20-alpine',
        '.',
        'dist',
        'npm install',
        'npm run build',
        'node dist/index.js'
    ),
    (
        'Python',
        'python:3.12-slim',
        '.',
        '.',
        'pip install -r requirements.txt',
        '',
        'python main.py'
    ),
    (
        'Go',
        'golang:1.23-alpine',
        '.',
        '.',
        'go mod download',
        'go build -o app ./cmd/...',
        './app'
    ),
    (
        'Static (Nginx)',
        'node:20-alpine',
        '.',
        'dist',
        'npm install',
        'npm run build',
        ''
    )
ON CONFLICT (name) DO NOTHING;
