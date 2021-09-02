import firestore from '../firestore';

const getSettings = () => {
  return firestore
    .collection('settings')
    .doc('default')
    .get()
    .then((docSnap) => {
      const data = docSnap.data();
      if (docSnap.exists && data) {
        return data;
      } else {
        throw new Error('Settings does not exists');
      }
    });
};

export default { getSettings };
