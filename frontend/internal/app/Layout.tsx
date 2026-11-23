import { Link, Outlet } from 'react-router'

// Shell shared by every page. Holds the navigation to the four catalog lists.
export function Layout() {
  return (
    <>
      <header>
        <nav>
          <Link to="/boardgames">Boardgames</Link>{' '}
          <Link to="/categories">Categories</Link>{' '}
          <Link to="/mechanisms">Mechanisms</Link>{' '}
          <Link to="/contributors">Contributors</Link>
        </nav>
      </header>
      <main>
        <Outlet />
      </main>
    </>
  )
}
