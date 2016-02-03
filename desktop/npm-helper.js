// TODO npm link helper for shared on windows
// Helper for cross platform npm run script commands
import path from 'path'
import {execSync} from 'child_process'

const [,,command, ...rest] = process.argv

const inject = info => {
  let temp = {
    ...info,
    env: {
      ...process.env,
      ...info.env
    }
  }

  if (info.nodeEnv) {
    temp.env.NODE_ENV = info.nodeEnv
  }

  if (info.nodePathDesktop) {
    temp.env.NODE_PATH = path.join(process.cwd(), 'node_modules')
  }

  if (rest.length && temp.shell) {
    temp.shell = temp.shell + ' ' + rest.join(' ')
  }

  return temp
}

const commands = {
  'start': () => {
    return {shell: 'npm run build-dev && npm run start-cold'}
  },
  'start-hot': () => {
    return {
      env: {HOT: 'true'},
      nodeEnv: 'development',
      shell: 'node client.js'
    }
  },
  'start-cold': () => {
    return {
      nodeEnv: 'development',
      shell: 'electron ./dist/main.bundle.js'
    }
  },
  'build-dev': () => {
    return {
      env: {NO_SERVER:'true', DEBUG: 'express:*'},
      nodeEnv: 'production',
      nodePathDesktop: true,
      shell: 'node server.js'}
  },
  'build-prod': () => {
    return {
      nodeEnv: 'production',
      nodePathDesktop: true,
      shell: 'webpack --config webpack.config.production.js --progress --profile --colors'
    }
  },
  'package': () => {
    return {
      nodeEnv: 'production',
      nodePathDesktop: true,
      shell: 'node package.js'
    }
  },
  'hot-server': () => {
    return {
      env: {HOT: 'true', DEBUG: 'express:*'},
      nodeEnv: 'production',
      nodePathDesktop: true,
      shell: 'node server.js'
    }
  },
  'inject-sourcemaps-prod': () => {
    return {shell: 'a(){ cp \'$1\'/* /Applications/Keybase.app/Contents/Resources/app/desktop/dist; };a'}
  },
  'inject-code-prod': () => {
    return {shell: 'npm run package; cp dist/* /Applications/Keybase.app/Contents/Resources/app/desktop/dist/'}
  },
  'start-prod': () => {
    return {shell: '/Applications/Keybase.app/Contents/MacOS/Electron'}
  },
  'electron-rebuild': () => {
    return {shell: './node_modules/.bin/electron-rebuild'}
  },
  'postinstall': () => {
    if (process.platform === 'windows') {
      fixupSymlinks()
    }
  }
}

function fixupSymlinks () {
  console.log('Fixing up shared symlinks')
  execSync('del shared && mklink /j shared ..\\shared', {cwd: path.join(process.cwd(), '.')})
  execSync('del shared && mklink /j shared ..\\shared', {cwd: path.join(process.cwd(), '..', 'react-native')})
}

const toRun = commands[command]

if (!toRun) {
  console.log('Unknown command: ', command)
  process.exit(1)
}

let info = toRun()

if (!info) {
  process.exit(0)
}

info = inject(info)

if (!info.shell) {
  process.exit(1)
}

console.log(`Calling: ${JSON.stringify(info, null, 2)}`)
const call = execSync(info.shell, {env: info.env, stdio: 'inherit'})
call.on('close', () => {})
