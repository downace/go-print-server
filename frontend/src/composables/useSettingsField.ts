import { readonly, shallowRef, watch } from "vue";

export default function useSettingsField<T = string>(
  initialValueGetter: () => T,
  updateValue: (newVal: T) => Promise<unknown>,
) {
  const field = shallowRef<T>(initialValueGetter());
  const error = shallowRef("");

  watch(field, async (value) => {
    try {
      await updateValue(value);
    } catch (e) {
      error.value = e instanceof Error ? e.message : (e as string);
    }
  });

  return {
    value: field,
    error: readonly(error),
  };
}
