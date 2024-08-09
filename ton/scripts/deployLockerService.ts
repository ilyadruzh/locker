import { toNano } from '@ton/core';
import { LockerService } from '../wrappers/LockerService';
import { compile, NetworkProvider } from '@ton/blueprint';

export async function run(provider: NetworkProvider) {
    const lockerService = provider.open(
        LockerService.createFromConfig(
            {
                id: Math.floor(Math.random() * 10000),
                counter: 0,
            },
            await compile('LockerService')
        )
    );

    await lockerService.sendDeploy(provider.sender(), toNano('0.05'));

    await provider.waitForDeploy(lockerService.address);

    console.log('ID', await lockerService.getID());
}
