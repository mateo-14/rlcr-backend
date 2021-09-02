import firestore from '../firestore';
import { client, dsRest } from '../ds';
import { Routes } from 'discord-api-types/v9';
import { RawUserData } from 'discord.js/typings/rawDataTypes';
import { User } from 'discord.js';

export const addOrUpdateUser = (id: string, token: string) => {
  return firestore
    .collection('users')
    .doc(id)
    .set({})
    .then(() => dsRest.put(Routes.guildMember(process.env.GUILD_ID!, id), { body: { access_token: token } }));
};
