import { Request, Response } from 'express';
import { dsRest } from '../ds';
import { URLSearchParams } from 'url';
import fetch from 'node-fetch';
import { Routes } from 'discord-api-types/v9';
import { addOrUpdateUser } from '../services/users';
import { RawUserData } from 'discord.js/typings/rawDataTypes';
import { generateToken } from '../services/tokens';

export const login = async (req: Request, res: Response) => {
  if (req.body.code) {
    const params = new URLSearchParams();
    params.append('client_id', process.env.CLIENT_ID!);
    params.append('client_secret', process.env.CLIENT_SECRET!);
    params.append('grant_type', 'authorization_code');
    params.append('code', req.body.code);
    params.append('redirect_uri', `${process.env.FRONTEND_URL}/ds_redirect`);

    try {
      const response = await fetch(`${process.env.DS_API_ENDPOINT}/oauth2/token`, {
        method: 'post',
        body: params,
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      });

      const oauth2 = await response.json();
      if (!response.ok) {
        throw new Error(`${response.status} ${JSON.stringify(oauth2)}`); // Will take you to the `catch` below
      }

      const user = (await dsRest.get(Routes.user(), {
        auth: false,
        headers: { Authorization: `Bearer ${oauth2.access_token}` },
      })) as RawUserData;

      await addOrUpdateUser(user.id, oauth2.access_token);
      const expireTime = parseInt(process.env.TOKEN_EXP_TIME!);
      const token = await generateToken(user.id, expireTime);
      res.cookie('token', token, { httpOnly: true, expires: new Date(Date.now() + expireTime) });
      res.json({ avatar: user.avatar, username: `${user.username}#${user.discriminator}`, id: user.id });
    } catch (err) {
      console.error(err);
      res.sendStatus(401);
    }
  }
};

export const logout = (_: Request, res: Response) => {
  res.cookie('token', '', { httpOnly: true, expires: new Date(0) });
  res.sendStatus(200);
};
