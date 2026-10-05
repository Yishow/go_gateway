/** Mirrors the existing canonical write-group destination capability in
 * write_group_repository_validation.go. Legacy connector kinds stay readable. */
export function supportsWriteGroupKind(kind: string): boolean {
  return kind === 'sqlite' || kind === 'postgres';
}
