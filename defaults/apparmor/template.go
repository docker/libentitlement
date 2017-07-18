package apparmor

/* See profile language http://wiki.apparmor.net/index.php/QuickProfileLanguage
 * This profile is the base template of Moby, customized to allow entitlements to configure
 * network, capabilities and file accesses
 */
const baseCustomTemplate = `
{{range $value := .Imports}}
{{$value}}
{{end}}

profile {{.Name}} flags=(attach_disconnected,mediate_deleted) {
{{range $value := .InnerImports}}
  {{$value}}
{{end}}

{{if .Network.Denied}}
  deny network
{{end}} {{else}}
  {{if .Network.AllowedProtocols}}
    {{range $value := .Network.AllowedProtocols}} network inet {{$value}}
  {{end}} {{else}}
    network,
  {{end}}
  {{if .Network.Raw.Denied}}
    deny network raw,
  {{end}}
{{end}}

{{range $value := .Capabilities.Allowed}} capabilty {{$value}}
{{range $value := .Capabilities.Denied}} deny capability {{$value}}

{{range $value := .Files.Denied}} deny {{$value}} rwamklx
{{range $value := .Files.ReadOnly}} deny {{$value}} wkal
{{range $value := .Files.NoExec}} deny {{$value}} x

  file,
  umount,

  deny @{PROC}/* w,   # deny write for all files directly in /proc (not in a subdir)
  # deny write to files not in /proc/<number>/** or /proc/sys/**
  deny @{PROC}/{[^1-9],[^1-9][^0-9],[^1-9s][^0-9y][^0-9s],[^1-9][^0-9][^0-9][^0-9]*}/** w,
  deny @{PROC}/sys/[^k]** w,  # deny /proc/sys except /proc/sys/k* (effectively /proc/sys/kernel)
  deny @{PROC}/sys/kernel/{?,??,[^s][^h][^m]**} w,  # deny everything except shm* in /proc/sys/kernel/
  deny @{PROC}/sysrq-trigger rwklx,
  deny @{PROC}/mem rwklx,
  deny @{PROC}/kmem rwklx,
  deny @{PROC}/kcore rwklx,

  deny mount,

  deny /sys/[^f]*/** wklx,
  deny /sys/f[^s]*/** wklx,
  deny /sys/fs/[^c]*/** wklx,
  deny /sys/fs/c[^g]*/** wklx,
  deny /sys/fs/cg[^r]*/** wklx,
  deny /sys/firmware/** rwklx,
  deny /sys/kernel/security/** rwklx,

{{if ge .Version 208095}}
  # suppress ptrace denials when using 'docker ps' or using 'ps' inside a container
  ptrace (trace,read) peer={{.Name}},
{{end}}
}
`