import { createClient } from '@supabase/supabase-js'

export const supabase = createClient(
  import.meta.env.VITE_APP_PROJECT_LINK,
  import.meta.env.VITE_APP_ANON_KEY
)
