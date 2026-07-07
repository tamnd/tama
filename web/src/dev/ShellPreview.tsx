import { MemoryRouter } from 'react-router-dom'
import { Shell } from '@/app/Shell'
import { Home } from '@/screens/Home/Home'
import { AuthProvider } from '@/state/auth'

// The gallery's shell specimen: the real Shell around the Home demo,
// mounted at /dev/gallery/shell so each resizable iframe gets a full
// document and the media queries respond to the frame, not the page.
export default function ShellPreview() {
  return (
    <AuthProvider>
      <MemoryRouter>
        <Shell>
          <Home />
        </Shell>
      </MemoryRouter>
    </AuthProvider>
  )
}
