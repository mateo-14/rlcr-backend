import firestore from '../firestore';
import { customAlphabet } from 'nanoid';
const nanoid = customAlphabet(process.env.NANOID_ALPHABET!, 30);

export const generateToken = (userID: string, expTime: number) => {
  const token = nanoid();
  return firestore
    .collection('tokens')
    .doc(token)
    .set({
      expDate: new Date(Date.now() + expTime),
      createdAt: new Date(),
      user: firestore.collection('users').doc(userID),
    })
    .then(() => token);
};

export const verify = (token: string) => {
  return firestore
    .collection('tokens')
    .doc(token)
    .get()
    .then(async (docSnap) => {
      let data;
      if (docSnap.exists && (data = docSnap.data())) {
        if (new Date() >= data.expDate.toDate()) {
          await docSnap.ref.delete();
          throw new Error('Token is expired');
        } else {
          return (data.user as FirebaseFirestore.DocumentReference).id;
        }
      } else {
        throw new Error('Token not exists');
      }
    });
};
