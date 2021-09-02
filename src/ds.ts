import { Client, Intents, MessageActionRow, MessageButton, MessageEmbed, User } from 'discord.js';

import { REST } from '@discordjs/rest';
import { Routes } from 'discord-api-types/v9';

const commands = [
  {
    name: 'pedidos',
    description: 'Muestra tus últimos pedidos o la información del pedido elegido',
  },
  {
    name: 'pedido',
    description: 'Muestra información del pedido',
    options: [{ type: 3, name: 'id', description: 'ID del pedido', required: true }],
  },
  /*  {
    name: 'pedir',
    description: 'Realiza un pedido',
    options: [
      {
        type: 4,
        name: 'tipo',
        description: 'Si vas a comprar o vender',
        required: true,
        choices: [
          { name: 'compra', value: 0 },
          { name: 'venta', value: 1 },
        ],
      },
      {
        type: 4,
        name: 'cantidad',
        description: 'Cantidad de créditos que queres comprar/vender (Para más opciones usa nuestra web)',
        required: true,
        choices: [
          { name: '100', value: 100 },
          { name: '200', value: 200 },
          { name: '300', value: 300 },
          { name: '400', value: 400 },
          { name: '500', value: 500 },
          { name: '600', value: 600 },
          { name: '700', value: 700 },
          { name: '800', value: 800 },
          { name: '900', value: 900 },
          { name: '1000', value: 1000 },
        ],
      },
      {
        type: 4,
        name: 'método',
        description: 'Método de pago',
        required: true,
        choices: [
          { name: 'Transferencia bancaria', value: 0 },
          { name: 'UALÁ', value: 1 },
          { name: 'MercadoPago', value: 2 },
        ],
      },
      {
        type: 3,
        name: 'perfil',
        description: 'Link al perfil de Steam o usuario de Epic Games',
        required: true,
      },
      {
        type: 3,
        name: 'cuenta',
        description: 'Cuenta de UALÁ/MercadoPago (Ingresar "-" si tu método de pago es Transferencia bancaria) ',
        required: true,
      },
      {
        type: 4,
        name: 'cbu',
        description: 'CBU/CVU (Ingresar 0 si tu método de pago es UALÁ o MercadoPago) ',
        required: true,
      },
      {
        type: 3,
        name: 'alias',
        description:
          'Alias de cuenta bancaria (Ingresar "-" si tu método de pago es UALÁ o MercadoPago, o si ingresaste CBU/CVU) ',
        required: true,
      },
      {
        type: 4,
        name: 'dni',
        description: 'Documento Nacional de Identidad (Ingresar 0 si tu método de pago es UALÁ o MercadoPago) ',
        required: true,
      },
    ],
  }, */
];

const rest = new REST({ version: '9' }).setToken(process.env.CLIENT_TOKEN!);
(async () => {
  try {
    console.log('Started refreshing application (/) commands.');
    await rest.put(Routes.applicationGuildCommands(process.env.CLIENT_ID!, process.env.GUILD_ID!), { body: commands });

    console.log('Successfully reloaded application (/) commands.');
  } catch (error) {
    console.error(error);
  }
})();

const client = new Client({ intents: [Intents.FLAGS.GUILDS] });

client.on('ready', (client) => {
  console.log(`Logged in as ${client.user.tag}!`);
});

client.on('interactionCreate', async (interaction) => {
  if (interaction.isCommand()) {
    if (interaction.commandName === 'pedidos') {
      const ordersURL = `${process.env.FRONTEND_URL!}/orders`;

      const embed = new MessageEmbed()
        .setColor('#8B5CF6')
        .setTimestamp(new Date())
        .setTitle('Tus pedidos')
        .setURL(`${process.env.FRONTEND_URL!}/orders`)
        .setDescription('Últimos # pedidos')
        .addField('Pedido #', `[Compra de # créditos a ARS$ #](${ordersURL})`);

      const row = new MessageActionRow().addComponents(
        new MessageButton().setLabel('Ver todos los pedidos').setStyle('LINK').setURL(ordersURL)
      );

      interaction.reply({ ephemeral: true, embeds: [embed], components: [row] });
    } else if (interaction.commandName === 'pedido') {
      const id = interaction.options.get('id')?.value;
      if (id) {
        const orderURL = `${process.env.FRONTEND_URL!}/orders/${id}`;

        const embed = new MessageEmbed()
          .setColor('#8B5CF6')
          .setTimestamp(new Date())
          .setTitle(`Pedido ${id}`)
          .setURL(orderURL)
          .setDescription('Info');
        const row = new MessageActionRow().addComponents(
          new MessageButton().setLabel('Ver pedido').setStyle('LINK').setURL(orderURL)
        );

        interaction.reply({ ephemeral: true, embeds: [embed], components: [row] });
      }
    } else if (interaction.commandName === 'pedir') {
    }
  }
});

const sendNewOrderMsg = (order: Order) => {
  client.users.fetch(order.userID).then((user) => {
    const orderURL = `${process.env.FRONTEND_URL!}/orders/${order.id}`;
    const embed = new MessageEmbed()
      .setColor('#8B5CF6')
      .setTitle(`Pedido realizado (${order.id})`)
      .setURL(orderURL)
      .setTimestamp(new Date()).setDescription(`**__Has realizado un pedido de ${
      order.mode === 0 ? 'compra' : 'venta'
    } de ${order.credits} créditos a ARS$ ${'precio'}__**\n
      • El pedido debe ser confirmado por un moderador. Una vez confirmado nos contactaremos por DM para realizar la transacción.
      • Si tenés algún problema o necesitas ayuda usa el comando **/ayuda** o contacta con un moderador en nuestro canal de discord.
      • Usa el comando **/pedidos** para ver la lista con los últimos pedidos.
      `);

    const row = new MessageActionRow().addComponents(
      new MessageButton().setLabel('Ver pedido').setStyle('LINK').setURL(orderURL)
    );

    user.send({ embeds: [embed], components: [row] });
  });
};
export { client, rest as dsRest, sendNewOrderMsg };