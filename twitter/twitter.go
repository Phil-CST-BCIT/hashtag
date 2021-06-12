import (
			"net"
			"io"
)	

// we want to keep tracking the state of a connection
var {
		cnx net.Conn 

		reader 	io.ReadCloser
	
	}

func dial(network, addr string) (net.Conn, error) {

	if cnx != nil {

		cnx.Close()
		
		cnx = nil
	
	}

	// establish a new connection and set timeout 3 sec
	conn, err := net.DialTimeout(network, addr, 3*time.Second)

	if err != nil {

		return nil, err
	
	}

	cnx = conn

	return conn, nil
}

func close_cnx() {
	
	if cnx != nil {
	
		cnx.Close()
	
	}

	if reader != nil {
	
		reader.Close()
	
	}

}