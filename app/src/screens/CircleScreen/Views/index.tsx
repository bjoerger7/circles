import React, { PureComponent } from 'react';
import { ScrollView, View, Text, Animated } from 'react-native';
import CircleMessageWindow from '../../../components/circle/MessageWindow';
import CircleNavbar from '../../../components/circle/Navbar';
import CircleFooter from '../../../components/circle/Footer';
import styles from './styles';

export interface Props {
  onPressBack(): void;
  footerHeight: Animated.Value;
  message: string;
  onMessageChange(newMessage: string): void;
}

class CircleScreenView extends PureComponent<Props> {
  render() {
    return (
      <View style={styles.container}>
        <CircleNavbar
          onPressBack={this.props.onPressBack}
          style={styles.headerContainer}
        />
        <CircleMessageWindow />
        <CircleFooter
          height={this.props.footerHeight}
          message={this.props.message}
          onMessageChange={this.props.onMessageChange}
        />
      </View>
    );
  }
}

export default CircleScreenView;
