//
// Sample server main program generated by sadl2java, with simple implementation added.
//
import org.eclipse.jetty.server.Server;
import org.glassfish.jersey.jetty.JettyHttpContainerFactory;
import org.glassfish.jersey.server.ResourceConfig;
import org.glassfish.hk2.utilities.binding.AbstractBinder;
import org.glassfish.jersey.jackson.JacksonFeature;
import javax.ws.rs.core.UriBuilder;
import java.io.IOException;
import java.net.URI;
import java.util.List;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Map;
import java.util.HashMap;
import swapi.model.*;
import swapi.server.*;

public class SwapiController implements Swapi {

	Map<String,Part> parts;
	Map<String,Film> films;
	
	public SwapiController() {
		parts = new HashMap<String,Part>();
		Part part1 = Part.builder().id("1").name("R2-D2").build();
		Part part2 = Part.builder().id("2").name("Han Solo").build();
		parts.put("1", part1);
		parts.put("2", part2);
		films = new HashMap<String,Film>();
		films.put("4", Film.builder().id("4").name("A New Hope").cast(Arrays.asList(part1, part2)).build());
		films.put("5", Film.builder().id("5").name("The Empire Strikes Back").build());
	}
	
	public GetPartResponse getPart(GetPartRequest req) {
		return GetPartResponse.builder().part(parts.get(req.getId())).build();
	}
	
	public GetFilmResponse getFilm(GetFilmRequest req) {
		return GetFilmResponse.builder().film(films.get(req.getId())).build();
	}
	
	public ListPartsResponse listParts(ListPartsRequest req) {
		ArrayList<Part> partValues = new ArrayList<Part>();
		for (Map.Entry<String,Part> kv : parts.entrySet()) {
			partValues.add(kv.getValue());
		}
		return ListPartsResponse.builder().parts(partValues).build();
	}
	
	public ListFilmsResponse listFilms(ListFilmsRequest req) {
		ArrayList<Film> filmValues = new ArrayList<Film>();
		for (Map.Entry<String,Film> kv : films.entrySet()) {
			filmValues.add(kv.getValue());
		}
		return ListFilmsResponse.builder().films(filmValues).build();
	}
	
}
