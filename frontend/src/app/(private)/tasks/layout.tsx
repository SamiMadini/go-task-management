export default async function Layout({
  children,
  sheet,
  modal,
}: {
  children: React.ReactNode
  sheet: React.ReactNode
  modal: React.ReactNode
}) {
  return (
    <>
      {children}
      {sheet}
      {modal}
    </>
  )
}
