import { Command } from 'commander';
import inquirer from 'inquirer';
import axios from 'axios';

const program = new Command();

const BASE_URL = process.env.USERS_SERVICE_URL || 'http://localhost:8081';

interface Driver {
  id: string;
  name: string;
  is_available: boolean;
}

interface Passenger {
  id: string;
  name: string;
  created_at: string;
  updated_at: string;
}

const api = axios.create({
  baseURL: BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

const handleError = (error: any) => {
  if (axios.isAxiosError(error)) {
    console.error('Error:', error.response?.data || error.message);
  } else {
    console.error('Error:', error.message);
  }
};

const displayIntro = () => {
  console.log('========================================');
  console.log('   MGL7361-Microservices implementation');
  console.log('   Users Service CLI');
  console.log('========================================');
  console.log('');
};

const setupDriverCommands = (program: Command) => {
  const driverCommand = program.command('drivers');
  
  driverCommand
    .command('create')
    .description('Create a new driver')
    .action(async () => {
      try {
        const answers = await inquirer.prompt([
          {
            type: 'input',
            name: 'name',
            message: 'Enter driver name:',
            validate: (input) => input.trim() ? true : 'Name is required'
          }
        ]);

        const response = await api.post('/drivers', {
          name: answers.name
        });

        console.log('Driver created successfully:');
        console.log(JSON.stringify(response.data, null, 2));
      } catch (error) {
        handleError(error);
      }
    });

  driverCommand
    .command('list')
    .description('List all drivers')
    .option('-a, --available', 'Show only available drivers')
    .action(async (options) => {
      try {
        const params = options.available ? { available: true } : {};
        const response = await api.get('/drivers', { params });
        
        console.log('Drivers:');
        console.log(JSON.stringify(response.data, null, 2));
      } catch (error) {
        handleError(error);
      }
    });

  driverCommand
    .command('update-status')
    .description('Update driver availability status')
    .action(async () => {
      try {
        const driversResponse = await api.get('/drivers');
        const drivers: Driver[] = driversResponse.data;

        const driverAnswer = await inquirer.prompt([
          {
            type: 'list',
            name: 'driverId',
            message: 'Select driver:',
            choices: drivers.map(driver => ({
              name: `${driver.name} (${driver.is_available ? 'Available' : 'Unavailable'})`,
              value: driver.id
            }))
          }
        ]);

        const statusAnswer = await inquirer.prompt([
          {
            type: 'confirm',
            name: 'isAvailable',
            message: 'Set as available?',
            default: true
          }
        ]);

        const response = await api.patch(`/drivers/${driverAnswer.driverId}/status`, {
          is_available: statusAnswer.isAvailable
        });

        console.log('Driver status updated successfully:');
        console.log(JSON.stringify(response.data, null, 2));
      } catch (error) {
        handleError(error);
      }
    });
};

const setupPassengerCommands = (program: Command) => {
  const passengerCommand = program.command('passengers');
  
  passengerCommand
    .command('create')
    .description('Create a new passenger')
    .action(async () => {
      try {
        const answers = await inquirer.prompt([
          {
            type: 'input',
            name: 'name',
            message: 'Enter passenger name:',
            validate: (input) => input.trim() ? true : 'Name is required'
          }
        ]);

        const response = await api.post('/passengers', {
          name: answers.name
        });

        console.log('Passenger created successfully:');
        console.log(JSON.stringify(response.data, null, 2));
      } catch (error) {
        handleError(error);
      }
    });

  passengerCommand
    .command('list')
    .description('List all passengers')
    .action(async () => {
      try {
        const response = await api.get('/passengers');
        
        console.log('Passengers:');
        console.log(JSON.stringify(response.data, null, 2));
      } catch (error) {
        handleError(error);
      }
    });

  passengerCommand
    .command('get')
    .description('Get passenger by ID')
    .action(async () => {
      try {
        const passengersResponse = await api.get('/passengers');
        const passengers: Passenger[] = passengersResponse.data;

        const answers = await inquirer.prompt([
          {
            type: 'list',
            name: 'passengerId',
            message: 'Select passenger:',
            choices: passengers.map(passenger => ({
              name: `${passenger.name} (${passenger.id})`,
              value: passenger.id
            }))
          }
        ]);

        const response = await api.get(`/passengers/${answers.passengerId}`);
        console.log('Passenger details:');
        console.log(JSON.stringify(response.data, null, 2));
      } catch (error) {
        handleError(error);
      }
    });

  passengerCommand
    .command('update')
    .description('Update passenger information')
    .action(async () => {
      try {
        const passengersResponse = await api.get('/passengers');
        const passengers: Passenger[] = passengersResponse.data;

        const passengerAnswer = await inquirer.prompt([
          {
            type: 'list',
            name: 'passengerId',
            message: 'Select passenger to update:',
            choices: passengers.map(passenger => ({
              name: `${passenger.name} (${passenger.id})`,
              value: passenger.id
            }))
          }
        ]);

        const nameAnswer = await inquirer.prompt([
          {
            type: 'input',
            name: 'name',
            message: 'Enter new name:',
            validate: (input) => input.trim() ? true : 'Name is required'
          }
        ]);

        const response = await api.put(`/passengers/${passengerAnswer.passengerId}`, {
          name: nameAnswer.name
        });

        console.log('Passenger updated successfully:');
        console.log(JSON.stringify(response.data, null, 2));
      } catch (error) {
        handleError(error);
      }
    });

  passengerCommand
    .command('delete')
    .description('Delete a passenger')
    .action(async () => {
      try {
        const passengersResponse = await api.get('/passengers');
        const passengers: Passenger[] = passengersResponse.data;

        const answers = await inquirer.prompt([
          {
            type: 'list',
            name: 'passengerId',
            message: 'Select passenger to delete:',
            choices: passengers.map(passenger => ({
              name: `${passenger.name} (${passenger.id})`,
              value: passenger.id
            }))
          },
          {
            type: 'confirm',
            name: 'confirm',
            message: 'Are you sure you want to delete this passenger?',
            default: false
          }
        ]);

        if (answers.confirm) {
          await api.delete(`/passengers/${answers.passengerId}`);
          console.log('Passenger deleted successfully');
        } else {
          console.log('Deletion cancelled');
        }
      } catch (error) {
        handleError(error);
      }
    });
};

const main = async () => {
  displayIntro();

  program
    .name('users-cli')
    .description('CLI for Users Service API')
    .version('1.0.0');

  setupDriverCommands(program);
  setupPassengerCommands(program);

  program
    .command('interactive')
    .description('Start interactive mode')
    .action(async () => {
      while (true) {
        const { action } = await inquirer.prompt([
          {
            type: 'list',
            name: 'action',
            message: 'What would you like to do?',
            choices: [
              { name: 'Manage Drivers', value: 'drivers' },
              { name: 'Manage Passengers', value: 'passengers' },
              { name: 'Exit', value: 'exit' }
            ]
          }
        ]);

        if (action === 'exit') {
          console.log('Goodbye!');
          process.exit(0);
        }

        if (action === 'drivers') {
          const { driverAction } = await inquirer.prompt([
            {
              type: 'list',
              name: 'driverAction',
              message: 'Driver Management:',
              choices: [
                { name: 'Create Driver', value: 'create' },
                { name: 'List Drivers', value: 'list' },
                { name: 'Update Driver Status', value: 'update-status' },
                { name: 'Back to Main Menu', value: 'back' }
              ]
            }
          ]);

          if (driverAction === 'back') continue;

          const driverCommand = program.commands.find(cmd => cmd.name() === 'drivers');
          if (driverCommand) {
            const subCommand = driverCommand.commands.find(cmd => cmd.name() === driverAction);
            if (subCommand) {
              await subCommand.parseAsync([], { from: 'user' });
            }
          }
        }

        if (action === 'passengers') {
          const { passengerAction } = await inquirer.prompt([
            {
              type: 'list',
              name: 'passengerAction',
              message: 'Passenger Management:',
              choices: [
                { name: 'Create Passenger', value: 'create' },
                { name: 'List Passengers', value: 'list' },
                { name: 'Get Passenger Details', value: 'get' },
                { name: 'Update Passenger', value: 'update' },
                { name: 'Delete Passenger', value: 'delete' },
                { name: 'Back to Main Menu', value: 'back' }
              ]
            }
          ]);

          if (passengerAction === 'back') continue;

          const passengerCommand = program.commands.find(cmd => cmd.name() === 'passengers');
          if (passengerCommand) {
            const subCommand = passengerCommand.commands.find(cmd => cmd.name() === passengerAction);
            if (subCommand) {
              await subCommand.parseAsync([], { from: 'user' });
            }
          }
        }
      }
    });

  // Mode Interactif
  const isDocker = process.env.IS_DOCKER || process.cwd() === '/app';
if (isDocker || process.argv.length <= 2) {
  console.log('Starting interactive mode...');
  program.commands.find(cmd => cmd.name() === 'interactive')?.parseAsync();
} else {
  program.parse();
}
};

main().catch(console.error);