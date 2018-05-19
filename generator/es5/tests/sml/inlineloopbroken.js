$module('inlineloopbroken', function () {
  var $static = this;
  $static.Value = $t.markpromising(function (props, value) {
    var $result;
    var $temp0;
    var $temp1;
    var v;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      localasyncloop: while (true) {
        switch ($current) {
          case 0:
            $current = 1;
            continue localasyncloop;

          case 1:
            $temp1 = value;
            $current = 2;
            continue localasyncloop;

          case 2:
            $promise.maybe($temp1.Next()).then(function ($result0) {
              $temp0 = $result0;
              $result = $temp0;
              $current = 3;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 3:
            v = $temp0.First;
            if ($temp0.Second.$wrapped) {
              $current = 4;
              continue localasyncloop;
            } else {
              $current = 5;
              continue localasyncloop;
            }
            break;

          case 4:
            $resolve($t.cast(v, $g.________testlib.basictypes.String, false));
            return;

          case 5:
            $resolve($t.fastbox('', $g.________testlib.basictypes.String));
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  });
  $static.Collector = $t.markpromising(function (props, chars) {
    var $result;
    var $temp0;
    var $temp1;
    var c;
    var final;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      localasyncloop: while (true) {
        switch ($current) {
          case 0:
            final = $t.fastbox('', $g.________testlib.basictypes.String);
            $current = 1;
            continue localasyncloop;

          case 1:
            $temp1 = chars;
            $current = 2;
            continue localasyncloop;

          case 2:
            $promise.maybe($temp1.Next()).then(function ($result0) {
              $temp0 = $result0;
              $result = $temp0;
              $current = 3;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 3:
            c = $temp0.First;
            if ($temp0.Second.$wrapped) {
              $current = 4;
              continue localasyncloop;
            } else {
              $current = 5;
              continue localasyncloop;
            }
            break;

          case 4:
            final = $g.________testlib.basictypes.String.$plus($t.cast(c, $g.________testlib.basictypes.String, false), final);
            $current = 2;
            continue localasyncloop;

          case 5:
            $resolve(final);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  });
  $static.TEST = $t.markpromising(function () {
    var $result;
    var characters;
    var result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      localasyncloop: while (true) {
        switch ($current) {
          case 0:
            characters = $t.fastbox('hello world', $g.________testlib.basictypes.String);
            $promise.maybe($g.inlineloopbroken.Collector($g.________testlib.basictypes.Mapping($t.any).Empty(), $g.________testlib.basictypes.MapStream($g.________testlib.basictypes.Integer, $g.________testlib.basictypes.String)($g.________testlib.basictypes.Integer.$exclusiverange($t.fastbox(0, $g.________testlib.basictypes.Integer), characters.Length()), $t.markpromising(function (index) {
              var $result;
              var $current = 0;
              var $continue = function ($resolve, $reject) {
                localasyncloop: while (true) {
                  switch ($current) {
                    case 0:
                      $promise.maybe($g.inlineloopbroken.Value($g.________testlib.basictypes.Mapping($t.any).Empty(), (function () {
                        var $current = 0;
                        var $continue = function ($yield, $yieldin, $reject, $done) {
                          while (true) {
                            switch ($current) {
                              case 0:
                                $yield(characters.$slice(index, $t.fastbox(index.$wrapped + 1, $g.________testlib.basictypes.Integer)));
                                $current = 1;
                                return;

                              case 1:
                                $done();
                                return;

                              default:
                                $done();
                                return;
                            }
                          }
                        };
                        return $generator.new($continue, false);
                      })())).then(function ($result0) {
                        $result = $result0;
                        $current = 1;
                        $continue($resolve, $reject);
                        return;
                      }).catch(function (err) {
                        $reject(err);
                        return;
                      });
                      return;

                    case 1:
                      $resolve($result);
                      return;

                    default:
                      $resolve();
                      return;
                  }
                }
              };
              return $promise.new($continue);
            })))).then(function ($result0) {
              $result = $result0;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            result = $result;
            $resolve($g.________testlib.basictypes.String.$equals(result, $t.fastbox('dlrow olleh', $g.________testlib.basictypes.String)));
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  });
});
