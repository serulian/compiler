$module('slice', function () {
  var $static = this;
  this.$struct('AnotherStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (AnotherInt) {
      var instance = new $static();
      instance.$data = {
      };
      instance.$lazycheck = false;
      instance.AnotherInt = AnotherInt;
      return $promise.resolve(instance);
    };
    $instance.Mapping = function () {
      var mappedData = {
      };
      mappedData['AnotherInt'] = this.AnotherInt;
      return $promise.resolve($t.nominalwrap(mappedData, $g.____testlib.basictypes.Mapping($t.any)));
    };
    $static.$equals = function (left, right) {
      if (left === right) {
        return $promise.resolve($t.nominalwrap(true, $g.____testlib.basictypes.Boolean));
      }
      var promises = [];
      promises.push($t.equals(left.$data['AnotherInt'], right.$data['AnotherInt'], $g.____testlib.basictypes.Integer));
      return Promise.all(promises).then(function (values) {
        for (var i = 0; i < values.length; i++) {
          if (!$t.unbox(values[i])) {
            return $t.nominalwrap(false, $g.____testlib.basictypes.Boolean);
          }
        }
        return $t.nominalwrap(true, $g.____testlib.basictypes.Boolean);
      });
    };
    Object.defineProperty($instance, 'AnotherInt', {
      get: function () {
        if (this.$lazycheck) {
          $t.ensurevalue(this.$data['AnotherInt'], $g.____testlib.basictypes.Integer, false, 'AnotherInt');
        }
        if (this.$data['AnotherInt'] != null) {
          return $t.box(this.$data['AnotherInt'], $g.____testlib.basictypes.Integer);
        }
        return this.$data['AnotherInt'];
      },
      set: function (val) {
        if (val != null) {
          this.$data['AnotherInt'] = $t.unbox(val);
          return;
        }
        this.$data['AnotherInt'] = val;
      },
    });
  });

  this.$struct('SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (Values) {
      var instance = new $static();
      instance.$data = {
      };
      instance.$lazycheck = false;
      instance.Values = Values;
      return $promise.resolve(instance);
    };
    $instance.Mapping = function () {
      var mappedData = {
      };
      mappedData['Values'] = this.Values;
      return $promise.resolve($t.nominalwrap(mappedData, $g.____testlib.basictypes.Mapping($t.any)));
    };
    $static.$equals = function (left, right) {
      if (left === right) {
        return $promise.resolve($t.nominalwrap(true, $g.____testlib.basictypes.Boolean));
      }
      var promises = [];
      promises.push($t.equals(left.$data['Values'], right.$data['Values'], $g.____testlib.basictypes.Slice($g.slice.AnotherStruct)));
      return Promise.all(promises).then(function (values) {
        for (var i = 0; i < values.length; i++) {
          if (!$t.unbox(values[i])) {
            return $t.nominalwrap(false, $g.____testlib.basictypes.Boolean);
          }
        }
        return $t.nominalwrap(true, $g.____testlib.basictypes.Boolean);
      });
    };
    Object.defineProperty($instance, 'Values', {
      get: function () {
        if (this.$lazycheck) {
          $t.ensurevalue(this.$data['Values'], $g.____testlib.basictypes.Slice($g.slice.AnotherStruct), false, 'Values');
        }
        if (this.$data['Values'] != null) {
          return $t.box(this.$data['Values'], $g.____testlib.basictypes.Slice($g.slice.AnotherStruct));
        }
        return this.$data['Values'];
      },
      set: function (val) {
        if (val != null) {
          this.$data['Values'] = $t.unbox(val);
          return;
        }
        this.$data['Values'] = val;
      },
    });
  });

  $static.TEST = function () {
    var correct;
    var jsonString;
    var parsed;
    var s;
    var values;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.slice.AnotherStruct.new($t.nominalwrap(1, $g.____testlib.basictypes.Integer)).then(function ($result0) {
              $temp0 = $result0;
              return $g.slice.AnotherStruct.new($t.nominalwrap(2, $g.____testlib.basictypes.Integer)).then(function ($result1) {
                $temp1 = $result1;
                return $g.slice.AnotherStruct.new($t.nominalwrap(3, $g.____testlib.basictypes.Integer)).then(function ($result2) {
                  $temp2 = $result2;
                  return $g.____testlib.basictypes.List($g.slice.AnotherStruct).forArray([($temp0, $temp0), ($temp1, $temp1), ($temp2, $temp2)]).then(function ($result3) {
                    $result = $result3;
                    $state.current = 1;
                    $callback($state);
                  });
                });
              });
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            values = $result;
            values.$slice($t.nominalwrap(0, $g.____testlib.basictypes.Integer), null).then(function ($result0) {
              return $g.slice.SomeStruct.new($result0).then(function ($result1) {
                $temp3 = $result1;
                $result = ($temp3, $temp3);
                $state.current = 2;
                $callback($state);
              });
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 2:
            s = $result;
            jsonString = $t.nominalwrap('{"Values":[{"AnotherInt":1},{"AnotherInt":2},{"AnotherInt":3}]}', $g.____testlib.basictypes.String);
            s.Stringify($g.____testlib.basictypes.JSON)().then(function ($result0) {
              return $g.____testlib.basictypes.String.$equals($result0, jsonString).then(function ($result1) {
                $result = $result1;
                $state.current = 3;
                $callback($state);
              });
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 3:
            correct = $result;
            $g.slice.SomeStruct.Parse($g.____testlib.basictypes.JSON)(jsonString).then(function ($result0) {
              $result = $result0;
              $state.current = 4;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 4:
            parsed = $result;
            $promise.resolve($t.nominalunwrap(correct)).then(function ($result1) {
              return ($promise.shortcircuit($result1, false) || s.Values.Length()).then(function ($result2) {
                return ($promise.shortcircuit($result1, false) || $g.____testlib.basictypes.Integer.$equals($result2, $t.nominalwrap(3, $g.____testlib.basictypes.Integer))).then(function ($result3) {
                  return $promise.resolve($result1 && $t.nominalunwrap($result3)).then(function ($result0) {
                    return ($promise.shortcircuit($result0, false) || s.Values.$index($t.nominalwrap(0, $g.____testlib.basictypes.Integer))).then(function ($result4) {
                      return ($promise.shortcircuit($result0, false) || $g.____testlib.basictypes.Integer.$equals($result4.AnotherInt, $t.nominalwrap(1, $g.____testlib.basictypes.Integer))).then(function ($result5) {
                        $result = $t.nominalwrap($result0 && $t.nominalunwrap($result5), $g.____testlib.basictypes.Boolean);
                        $state.current = 5;
                        $callback($state);
                      });
                    });
                  });
                });
              });
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 5:
            $state.resolve($result);
            return;

          default:
            $state.current = -1;
            return;
        }
      }
    });
    return $promise.build($state);
  };
});
