// The tab view: one entry, its state kept with the tab (across reloads, desks and splits).
import { Show } from 'solid-js'
import { useIde, useViewState } from '@antennae-space/plugin-sdk/solid'
import { Button, Field, TextArea, View } from '@antennae-space/plugin-sdk/ui'
import type { Entry } from './types'

export function EntryView() {
  const ide = useIde()
  // The entry is the state the tool opened the tab with; the note is added by this view.
  const [state, setState] = useViewState<Entry & { note?: string }>({ name: '', path: '', dir: false, size: 0 })
  return (
    <View>
      <div class="example-entry">
        <h2>{state().name}</h2>
        <p class="ide-muted">{state().dir ? ide.t('A folder of the project') : ide.t('A file of {n} bytes', { n: state().size })}</p>
        <Show when={!state().dir}>
          <Button icon="file" onClick={() => ide.openFile(state().path)}>{ide.t('Open the file')}</Button>
        </Show>
        <Field label={ide.t('A note, kept with the tab')}>
          <TextArea rows={4} value={state().note ?? ''} onInput={(note) => setState({ ...state(), note })} />
        </Field>
      </div>
    </View>
  )
}
