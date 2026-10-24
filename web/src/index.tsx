// The page module of the plugin: what it registers in the IDE, by the ids of plugin.json. The
// IDE runs it in a sandboxed frame (no network, no storage, no access to the page) and gives it
// `ide`, the sandbox API; the SDK's preset builds it into web/index.js.
import { view } from '@antennae-space/plugin-sdk/solid'
import { Tool } from './Tool'
import { EntryView } from './Entry'
import type { Summary } from './types'
import './example.css'

ide.register({
  views: {
    // The tool of the rail (contributes.tools) and the tab view (contributes.tabViews).
    'example': view(Tool),
    entry: view(EntryView),
  },
  commands: {
    // contributes.commands: run in a hidden frame, from Search Everywhere or its shortcut.
    summary: async (ide) => {
      try {
        const s = await ide.rpc<Summary>('example.summary')
        ide.toast(ide.t('{folders} folders and {files} files at the root of the project', { ...s }), 'ok')
      } catch (e) {
        ide.toast(e instanceof Error ? e.message : String(e), 'error')
      }
    },
  },
})
