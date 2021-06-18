$params = @{
  Name = "wireguard_exporter"
  BinaryPathName = '"C:\Program Files\WireGuard\wireguard_exporter.exe" --log.format logger:eventlog?name=wireguard_exporter --metrics.addr 127.0.0.1:9586'
  DependsOn = "WireGuardManager"
  DisplayName = "WireGuard Metrics"
  StartupType = "Automatic"
  Description = "Exports prometheus metrics about WireGuard."
}
New-Service @params

function Set-ServiceRecovery{
  [alias('Set-Recovery')]
  param
  (
      [string] [Parameter(Mandatory=$true)] $ServiceName,
      [string] $action1 = "restart",
      [int] $time1 =  30000, # in miliseconds
      [string] $action2 = "restart",
      [int] $time2 =  30000, # in miliseconds
      [string] $actionLast = "restart",
      [int] $timeLast = 30000, # in miliseconds
      [int] $resetCounter = 4000 # in seconds
  )
  $services = Get-CimInstance -ClassName 'Win32_Service' | Where-Object {$_.Name -imatch $ServiceName}
  $action = $action1+"/"+$time1+"/"+$action2+"/"+$time2+"/"+$actionLast+"/"+$timeLast
  foreach ($service in $services){
      # https://technet.microsoft.com/en-us/library/cc742019.aspx
      $output = sc.exe failure $($service.Name) actions= $action reset= $resetCounter
  }
}
Set-ServiceRecovery -ServiceName "wireguard_exporter"