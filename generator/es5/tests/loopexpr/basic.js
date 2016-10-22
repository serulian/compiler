$module('basic', function () {
  var $static = this;
  $static.SomeGenerator = function () {
    var $current = 0;
    var $continue = function ($yield, $yieldin, $reject, $done) {
      while (true) {
        switch ($current) {
          case 0:
            $yield($t.fastbox(1, $g.____testlib.basictypes.Integer));
            $current = 1;
            return;

          case 1:
            $yield($t.fastbox(2, $g.____testlib.basictypes.Integer));
            $current = 2;
            return;

          case 2:
            $yield($t.fastbox(3, $g.____testlib.basictypes.Integer));
            $current = 3;
            return;

          default:
            $done();
            return;
        }
      }
    };
    return $generator.new($continue);
  };
  $static.TEST = function () {
    var $result;
    var $temp0;
    var $temp1;
    var counter;
    var entry;
    var s;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.basic.SomeGenerator().then(function ($result1) {
              return $g.____testlib.basictypes.MapStream($g.____testlib.basictypes.Integer, $g.____testlib.basictypes.Integer)($result1, function (s) {
                var $result;
                var $current = 0;
                var $continue = function ($resolve, $reject) {
                  while (true) {
                    switch ($current) {
                      case 0:
                        $g.____testlib.basictypes.Integer.$plus(s, $t.fastbox(1, $g.____testlib.basictypes.Integer)).then(function ($result0) {
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
              }).then(function ($result0) {
                $result = $result0;
                $current = 1;
                $continue($resolve, $reject);
                return;
              });
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            s = $result;
            counter = $t.fastbox(0, $g.____testlib.basictypes.Integer);
            $current = 2;
            continue;

          case 2:
            $temp1 = s;
            $current = 3;
            continue;

          case 3:
            $temp1.Next().then(function ($result0) {
              $temp0 = $result0;
              $result = $temp0;
              $current = 4;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 4:
            entry = $temp0.First;
            if ($temp0.Second.$wrapped) {
              $current = 5;
              continue;
            } else {
              $current = 7;
              continue;
            }
            break;

          case 5:
            $g.____testlib.basictypes.Integer.$plus(counter, entry).then(function ($result0) {
              counter = $result0;
              $result = counter;
              $current = 6;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 6:
            $current = 3;
            continue;

          case 7:
            $g.____testlib.basictypes.Integer.$equals(counter, $t.fastbox(9, $g.____testlib.basictypes.Integer)).then(function ($result0) {
              $result = $result0;
              $current = 8;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 8:
            $resolve($result);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
});
