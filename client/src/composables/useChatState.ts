import { ref } from 'vue';
import { useRoute } from 'vue-router';

export interface IUser {
  id: string;
  username: string;
  avatar: string;
  doneLoading?: boolean;
}

export function isUser(object: Record<any, any>): object is IUser {
  return 'username' in object;
}

export interface IRoom {
  id: string;
  name: string;
  doneLoading?: boolean;
}

export function isRoom(object: Record<any, any>): object is IRoom {
  return 'name' in object;
}

export enum MessageType {
  NEW_MESSAGE,
  ME_CHANGED_USERNAME,
  OTHER_CHANGED_USERNAME,
  OTHER_LEFT,
  ME_JOINED,
  OTHER_JOINED
}

export interface IMessage {
  id: string;
  message: string;
  timestamp: number;
  user: IUser;
  type: MessageType;
}

const me = ref<IUser>({
  id: '',
  username: '',
  // avatar:
  //   'https://avataaars.io/?avatarStyle=Circle&topType=WinterHat1&accessoriesType=Kurt&hatColor=Red&facialHairType=MoustacheFancy&facialHairColor=BrownDark&clotheType=BlazerShirt&eyeType=Default&eyebrowType=Default&mouthType=Default&skinColor=Light'
  avatar: 'https://picsum.photos/56/56'
});
const messageInput = ref('');
const messages = ref<IMessage[]>([]);
const currentRoom = ref<IRoom | IUser>();
const rooms = ref<IRoom[]>([]);
const users = ref<IUser[]>([]);

export default function useChatState() {
  const route = useRoute();
  // const currentRoom = computed(() => {
  //   const isfound = rooms.value.find((room) => room.id === route.query.roomId);
  //   if (isfound) {
  //     return isfound;
  //   }
  //   return users.value.find((user) => user.id === route.query.roomId);
  // });

  return {
    rooms,
    users,
    me,
    messageInput,
    messages,
    currentRoom
  };
}
