import firestore from '../firestore';
import { User } from 'discord.js';

export const addOrUpdate = (user: User) => {
  return firestore
    .collection('users')
    .doc(user.id)
    .set({ username: user.username, avatar: user.avatar, discriminator: user.discriminator }, { merge: true });
};

export const getData = (id: string) => {
  return firestore
    .collection('users')
    .doc(id)
    .get()
    .then((docSnap) => {
      if (docSnap.exists) {
        return docSnap.data();
      } else {
        throw new Error('User does not exists');
      }
    });
};

export const getAll = () => {
  return firestore
    .collection('users')
    .get()
    .then((querySnap) => querySnap.docs.map((docSnap) => ({ ...docSnap.data(), id: docSnap.id })));
};
