program Hypervisor{
    licence "MIT"
    author  "Ali Rankine"
    description "Interacts with system device and spawns proceses"


    import std:id 
    
    export type SpawnCommand {
        programName:   string,
    }

    const debug: bool = true

    let pending = map[id.ID]VMID @capacity(32)

    attach device.switch
    attach device.system

    @comptime if debug{
        attach device.terminal
    }


    // signal 0 reserved for startup, handlers automatically registered if startup handler not defined
    export signal spawn(SpawnCommand) = 1
    export signal close(VMID) = 2
    export signal system.spawned(requestID: id.ID, process: Result<system.SpawnedProcess>) = 3

    handle spawn(command: SpawnCommand){
        let requestID = id.new()
        pending[requestID] = @caller
        system.spawn <- system.SpawnRequest{ID:requestID,Name: command.programName}
        debugPrint("process spawn requested")
        
    }

    // namespaced signal can only come from a device named 'system'
    handle system.spawned(requestID: id.ID, process: Result<system.SpawnedProcess>){
        let caller = pending[requestID]
        defer delete(pending, requestID)
        match process{
            Ok(spawned): {
                debugPrint("process spawn completed")
                switch.send[caller] <-  Result<VMID>.Ok(spawned.VMID)
            }
            Err(error): {
                debugPrint("process spawn errored")
                switch.send[caller] <-  Result<VMID>.Err(error) 
            }
        }
    }

    handle close(vmID: VMID){
        system.kill <- vmID
        switch.send[@caller] <- Result.Ok() // void typed result, can be an error, or empty (<void> type parameter is elided)
    }

    // if debug is false, the function body is empty, and therefore calls will be elided.
    fn debugPrint(text: string){
        @comptime if debug{
            terminal.writeOut <- text
        }
    }

    fn sum(nums: [u8]): u8{
        let acc: u8 = 0
        for i in len(nums){
            acc+=nums[i]
        }
        return acc
    }

    fn arbitraryAdd(): u8{
        @asm(
            push 0x0F
            push 0x1F
            add
        )
        return @stack(1)
    }

    fn swap(a: u8, b: u8): u8,u8{
        @asm(
            swp
        )
        return @stack(2)
    }

    fn greet(name: string): string{
        return switch name{
            case "Ali":
                "Hello friend"
            default:
                "Hello stranger"
        } 
    }
}