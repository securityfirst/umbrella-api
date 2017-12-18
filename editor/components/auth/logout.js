import { Component, PropTypes } from 'react'
import { connect } from 'react-redux'
import { withRouter } from 'react-router'
import * as authActionCreators from '../../actions/auth'

class LogoutPage extends Component {

  componentWillMount() {
    this.props.dispatch(authActionCreators.logoutUser())
    this.props.router.replace('/login')
  }

  render() {
    return null
  }
}

LogoutPage.propTypes = {
  dispatch: PropTypes.func.isRequired,
  router: PropTypes.object.isRequired
}

export default connect()(withRouter(LogoutPage))