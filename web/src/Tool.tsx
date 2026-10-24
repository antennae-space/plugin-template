// The tool: the entries at the root of the project, from the backend, as a list or a table.
// Enter or a double click opens an entry in a tab of the plugin; the backend's event
// example.changed reloads every window of the project.
import { createSignal, Show } from 'solid-js'
import { useEvent, useIde, useRpc } from '@antennae-space/plugin-sdk/solid'
import { EmptyState, IconButton, List, Message, Spinner, Table, Tabs, View } from '@antennae-space/plugin-sdk/ui'
import type { Entry } from './types'

export function Tool() {
  const ide = useIde()
  const entries = useRpc<Entry[]>('example.entries')
  const [mode, setMode] = createSignal<'list' | 'table'>('list')
  const [selected, setSelected] = createSignal<string>()
  useEvent('example.changed', () => void entries.reload())
  const open = (e: Entry) => ide.openView('entry', e, e.name)
  const size = (n: number) => (n < 1024 ? `${n} B` : `${(n / 1024).toFixed(1)} KB`)

  return (
    <View
      toolbar={
        <>
          <strong class="example-title">{ide.t('Project root')}</strong>
          <Show when={entries.loading()}><Spinner title={ide.t('Loading')} /></Show>
          <IconButton icon="refresh" title={ide.t('Reload in every window')} onClick={() => void ide.rpc('example.refresh')} />
        </>
      }
    >
      <Tabs tabs={[{ id: 'list', label: ide.t('List') }, { id: 'table', label: ide.t('Table') }]} active={mode()} onSelect={setMode} />
      <Show when={entries.error()}>{(err) => <Message kind="error">{err()}</Message>}</Show>
      <Show when={entries.value()?.length} fallback={<Show when={!entries.loading() && !entries.error()}><EmptyState text={ide.t('The project is empty')} /></Show>}>
        <Show
          when={mode() === 'list'}
          fallback={
            <Table
              label={ide.t('Entries')}
              columns={[
                { id: 'name', label: ide.t('Name'), value: (e: Entry) => e.name },
                { id: 'size', label: ide.t('Size'), value: (e: Entry) => (e.dir ? '' : size(e.size)), align: 'right' },
              ]}
              rows={entries.value() ?? []}
              onOpen={open}
            />
          }
        >
          <List
            label={ide.t('Entries')}
            items={(entries.value() ?? []).map((e) => ({ id: e.path, label: e.name, icon: e.dir ? ('folder' as const) : ('file' as const), detail: e.dir ? '' : size(e.size) }))}
            selected={selected()}
            onSelect={setSelected}
            onOpen={(it) => open(entries.value()!.find((e) => e.path === it.id)!)}
          />
        </Show>
      </Show>
    </View>
  )
}
