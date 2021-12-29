import { Client, Intents, MessageActionRow, MessageButton, MessageEmbed, User } from 'discord.js';

import { REST } from '@discordjs/rest';
import { Routes } from 'discord-api-types/v9';
import { URLSearchParams } from 'url';
import fetch from 'node-fetch';

const rest = new REST({ version: '9' }).setToken(process.env.CLIENT_TOKEN!);
const client = new Client({ intents: [Intents.FLAGS.GUILDS, Intents.FLAGS.GUILD_MEMBERS] });

client.on('ready', async (client) => {
  console.log(`Logged in as ${client.user.tag}!`);
});

const sendNewOrderMsg = (order: Order) => {
  client.users.fetch(order.userID).then(async (user) => {
    const orderURL = `${process.env.FRONTEND_URL!}/orders/${order.id}`;
    const embed = new MessageEmbed()
      .setColor('#8B5CF6')
      .setTitle(`Pedido realizado (${order.id})`)
      .setURL(orderURL)
      .setTimestamp(new Date()).setDescription(
        `**__Has realizado un pedido de ${order.mode === 0 ? 'compra' : 'venta'} 
        de ${order.credits} créditos a ARS$ ${order.price}__**\n
        **\*\*ESTO ES UNA DEMO\*\***
      `);

    const row = new MessageActionRow().addComponents(
      new MessageButton().setLabel('Ver pedido').setStyle('LINK').setURL(orderURL)
    );

    try {
      await user.send({ embeds: [embed], components: [row] });

      // Notify moderators
      const adminOrderUrl = `${process.env.FRONTEND_URL!}/admin/orders?id=${order.id}`;
      const moderatorEmbed = new MessageEmbed()
        .setColor('#8B5CF6')
        .setTitle(`Nuevo pedido (${order.id})`)
        .setURL(adminOrderUrl)
        .setTimestamp(new Date()).setDescription(`**El usuario ${user.username}#${user.discriminator} (${
        user.id
      }) ha realizado un pedido de ${order.mode === 0 ? 'compra' : 'venta'} de ${order.credits} créditos a ARS$ ${
        order.price
      }**
    `);

      const moderatorRow = new MessageActionRow().addComponents(
        new MessageButton().setLabel('Ver pedido').setStyle('LINK').setURL(adminOrderUrl)
      );

      const guild = client.guilds.cache.get(process.env.GUILD_ID!);
      await guild?.members.fetch();
      const members = guild?.members.cache.filter((member) =>
        member.roles.cache.has(process.env.DS_MODERATOR_ROLE_ID!)
      );
      members?.forEach(({ user }) => {
        if (!user.bot) {
          user.send({ embeds: [moderatorEmbed], components: [moderatorRow] });
        }
      });
    } catch {}
  });
};

const oauth2ByCode = (code: string): Promise<string> => {
  const params = new URLSearchParams();
  params.append('client_id', process.env.CLIENT_ID!);
  params.append('client_secret', process.env.CLIENT_SECRET!);
  params.append('grant_type', 'authorization_code');
  params.append('code', code);
  params.append('redirect_uri', `${process.env.FRONTEND_URL}/ds_redirect`);
  return fetch(`${process.env.DS_API_ENDPOINT}/oauth2/token`, {
    method: 'post',
    body: params,
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
  })
    .then((response) => {
      if (response.ok) return response.json();
      else throw new Error(`${response.status} ${response.statusText}`);
    })
    .then((response) => response.access_token as string);
};

const getUserByToken = (token: string): Promise<User> => {
  return rest
    .get(Routes.user(), {
      auth: false,
      headers: { Authorization: `Bearer ${token}` },
    })
    .then((user) => user as User);
};

const addUserToGuild = (id: string, token: string) => {
  return rest.put(Routes.guildMember(process.env.GUILD_ID!, id), { body: { access_token: token } });
};

export { client, sendNewOrderMsg, oauth2ByCode, addUserToGuild, getUserByToken };
