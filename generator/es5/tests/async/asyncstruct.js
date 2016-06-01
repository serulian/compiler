$module('asyncstruct', function () {
  var $static = this;
  this.$struct('SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (Foo, Bar) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        Foo: Foo,
        Bar: Bar,
      };
      return $promise.resolve(instance);
    };
    $static.$box = function (data) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = data;
      instance.$lazychecked = {
      };
      Object.defineProperty(instance, 'Foo', {
        get: function () {
          if (this.$lazychecked['Foo']) {
            $t.ensurevalue(this[BOXED_DATA_PROPERTY]['Foo'], $g.____testlib.basictypes.Integer, false, 'Foo');
            this.$lazychecked['Foo'] = true;
          }
          return $t.box(this[BOXED_DATA_PROPERTY]['Foo'], $g.____testlib.basictypes.Integer);
        },
      });
      Object.defineProperty(instance, 'Bar', {
        get: function () {
          if (this.$lazychecked['Bar']) {
            $t.ensurevalue(this[BOXED_DATA_PROPERTY]['Bar'], $g.____testlib.basictypes.Integer, false, 'Bar');
            this.$lazychecked['Bar'] = true;
          }
          return $t.box(this[BOXED_DATA_PROPERTY]['Bar'], $g.____testlib.basictypes.Integer);
        },
      });
      instance.Mapping = function () {
        var mapped = {
        };
        mapped['Foo'] = this.Foo;
        mapped['Bar'] = this.Bar;
        return $promise.resolve($t.box(mapped, $g.____testlib.basictypes.Mapping($t.any)));
      };
      return instance;
    };
    $instance.Mapping = function () {
      return $promise.resolve($t.box(this[BOXED_DATA_PROPERTY], $g.____testlib.basictypes.Mapping($t.any)));
    };
    $static.$equals = function (left, right) {
      if (left === right) {
        return $promise.resolve($t.box(true, $g.____testlib.basictypes.Boolean));
      }
      var promises = [];
      promises.push($t.equals(left[BOXED_DATA_PROPERTY]['Foo'], right[BOXED_DATA_PROPERTY]['Foo'], $g.____testlib.basictypes.Integer));
      promises.push($t.equals(left[BOXED_DATA_PROPERTY]['Bar'], right[BOXED_DATA_PROPERTY]['Bar'], $g.____testlib.basictypes.Integer));
      return Promise.all(promises).then(function (values) {
        for (var i = 0; i < values.length; i++) {
          if (!$t.unbox(values[i])) {
            return values[i];
          }
        }
        return $t.box(true, $g.____testlib.basictypes.Boolean);
      });
    };
    Object.defineProperty($instance, 'Foo', {
      get: function () {
        return this[BOXED_DATA_PROPERTY]['Foo'];
      },
      set: function (value) {
        this[BOXED_DATA_PROPERTY]['Foo'] = value;
      },
    });
    Object.defineProperty($instance, 'Bar', {
      get: function () {
        return this[BOXED_DATA_PROPERTY]['Bar'];
      },
      set: function (value) {
        this[BOXED_DATA_PROPERTY]['Bar'] = value;
      },
    });
  });

  $static.DoSomethingAsync = $t.workerwrap('dd7aa26ec2db6858f13a458fe077b54060448d2e8a2998b7445a87f974ab256b', function (s) {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            s.Foo.String().then(function ($result1) {
              return s.Bar.String().then(function ($result2) {
                return $g.____testlib.basictypes.String.$plus($result1, $result2).then(function ($result0) {
                  $result = $result0;
                  $current = 1;
                  $continue($resolve, $reject);
                  return;
                });
              });
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
  });
  $static.TEST = function () {
    var vle;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.asyncstruct.SomeStruct.new($t.box(1, $g.____testlib.basictypes.Integer), $t.box(2, $g.____testlib.basictypes.Integer)).then(function ($result1) {
              $temp0 = $result1;
              return $promise.translate($g.asyncstruct.DoSomethingAsync(($temp0, $temp0))).then(function ($result0) {
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
            vle = $result;
            $g.____testlib.basictypes.String.$equals(vle, $t.box("12", $g.____testlib.basictypes.String)).then(function ($result0) {
              $result = $result0;
              $current = 2;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 2:
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
